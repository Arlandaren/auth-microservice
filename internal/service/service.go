package service

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"service/internal/repository"
	"service/internal/shared/storage/dto"
	"service/internal/shared/utils"
	"service/internal/transport/grpc"
	pb "service/pkg/grpc/auth_v1"
	"strconv"
	"strings"
	"time"
)

const Issuer = "auth-microservice"

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	log.Println("NewService")
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	existingUser, err := s.repo.GetUserByName(req.Name)
	if existingUser != nil {
		return nil, fmt.Errorf("user with given creds is existing: %w", err)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("registering user")
		} else {
			return nil, err
		}

	}
	hash, err := utils.GenerateHashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	log.Println("hash", hash)
	user := dto.User{
		Name:     req.Name,
		Password: hash,
		Role:     dto.UserRole,
	}

	err = s.repo.NewUser(&user)
	if err != nil {
		return nil, err
	}
	log.Println("user created: ", user)
	tokenString, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	log.Println("token: ", tokenString)

	return &pb.RegisterResponse{Token: tokenString}, nil
}

func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := s.repo.GetUserByName(req.Name)
	if err != nil {
		return nil, errors.New("invalid name or password")
	}

	err = utils.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, errors.New("invalid name or password")
	}

	tokenString, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}

func (s *Service) RegisterAdmin(ctx context.Context, req *pb.RegisterAdminRequest, role string) (*pb.RegisterAdminResponse, error) {
	if strings.TrimSpace(role) != strings.TrimSpace(dto.SupremeRole) {
		return nil, errors.New("access denied, you are not admin")
	}

	existingUser, err := s.repo.GetUserByName(req.Name)
	if existingUser != nil {
		return nil, fmt.Errorf("user with given creds is existing: %w", err)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("registering user")
		} else {
			return nil, err
		}

	}
	hash, err := utils.GenerateHashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	log.Println("hash", hash)
	user := dto.User{
		Name:     req.Name,
		Password: hash,
		Role:     req.Role,
	}

	err = s.repo.NewUser(&user)
	if err != nil {
		return nil, err
	}
	log.Println("admin-user created: ", user)

	return &pb.RegisterAdminResponse{Id: int64(user.ID)}, nil
}

func (s *Service) OIDCToken(ctx context.Context, req *pb.OIDCTokenRequest) (*pb.OIDCTokenResponse, error) {
	// Проверяем ClientID на существование в DB
	err := s.repo.CheckGetClientID(req.ClientId)
	if err != nil {
		return nil, err
	}

	// Проверяем пользователя на нахождение в DB
	user, err := s.repo.GetUserByName(req.Name)
	if err != nil {
		return nil, errors.New("name does not exist")
	}
	err = utils.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, errors.New("invalid name or password")
	}

	// Валидация и подставление значений
	resultScopes, err := utils.AddScopes(user, req.Scopes, utils.TokenID)
	if err != nil {
		return nil, err
	}
	// Генерируем idToken
	idToken, err := utils.GenerateTokenOIDC(strconv.Itoa(user.ID), req.ClientId,
		Issuer, resultScopes, utils.PrivateKeyAccessToken)
	if err != nil {
		return nil, fmt.Errorf("не удалось сгенерировать id_token, error: %q", err)
	}

	// Генерируем authCode
	authCode, err := utils.GenerateAuthCode(user.ID)
	if err != nil {
		return nil, err
	}

	authCodeDTO := dto.AuthCodeOIDC{
		AuthCode:            authCode,
		UserID:              user.ID,
		ClientID:            req.ClientId,
		RedirectURI:         req.RedirectUri,
		State:               req.State,
		CodeChallengeMethod: req.CodeChallengeMethod,
		CodeChallenge:       req.CodeChallenge,
		ExpiresAt:           time.Now().Add(8 * 60 * time.Second),
	}
	if err := s.repo.NewAuthCode(&authCodeDTO); err != nil {
		return nil, err
	}

	response := &pb.OIDCTokenResponse{
		State:    req.State,
		IdToken:  idToken,
		AuthCode: authCode,
	}
	return response, nil
}

func (s *Service) OIDCExchange(ctx context.Context, req *pb.OIDCExchangeRequest) (*pb.OIDCExchangeResponse, error) {
	// Запросы в базу
	authCodeOIDC, err := s.repo.GetAuthCodeFromClientID(req.ClientId)
	if err != nil {
		return nil, err
	}
	clientData, err := s.repo.GetClientIDandClientSecret(req.ClientId)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetUserById(authCodeOIDC.UserID)
	if err != nil {
		return nil, err
	}

	// Проверки перед выдачей accessToken
	if req.ClientSecret != clientData.ClientSecret {
		return nil, errors.New("неверный ClientSecret")
	}
	if req.Code != authCodeOIDC.AuthCode {
		return nil, errors.New("неверный AuthCode")
	}
	if req.RedirectUri != authCodeOIDC.RedirectURI {
		return nil, errors.New("uri при получении AuthCode и при получении AccessToken не совпадают")
	}
	if time.Now().After(authCodeOIDC.ExpiresAt) {
		err := s.repo.DeleteAuthCodeFromClientID(req.ClientId)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("время действия AuthCode истекло")
	}
	if err := utils.VerifyCodeVerifier(req.CodeVerifier, authCodeOIDC.CodeChallengeMethod, authCodeOIDC.CodeChallenge); err != nil {
		return nil, err
	}

	// Валидация и подставление значений
	resultScopes, err := utils.AddScopes(user, authCodeOIDC.Scopes, utils.AccessToken)
	if err != nil {
		return nil, err
	}
	accessToken, err := utils.GenerateTokenOIDC(strconv.Itoa(authCodeOIDC.UserID), authCodeOIDC.ClientID,
		Issuer, resultScopes, utils.AccessToken)

	TokenDTO := dto.AccessToken{
		Token:     accessToken,
		ClientID:  authCodeOIDC.ClientID,
		UserID:    authCodeOIDC.UserID,
		TokenType: "bearer",
		ExpiresAt: time.Now().Add(60 * 60 * time.Second),
	}

	// Добавление AccessToken в базу
	if err := s.repo.NewAccessToken(&TokenDTO); err != nil {
		return nil, err
	}

	// Отправляем токен клиенту
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := grpc.SendClientRequest(ctx, authCodeOIDC.RedirectURI, &TokenDTO); err != nil {
		return nil, err
	}

	return &pb.OIDCExchangeResponse{}, nil
}
