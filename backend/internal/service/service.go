package service

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"service/internal/shared/config"
	"strings"

	log "github.com/sirupsen/logrus"
	"service/internal/repository"
	"service/internal/shared/storage/dto"
	"service/internal/shared/utils"
	pb "service/pkg/grpc/auth_v1"
)

type Service struct {
	Repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	log.Println("NewService")
	return &Service{
		Repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	forbiddenRoles := []string{"Supreme", "Client_Supreme"}
	for _, forbiddenRole := range forbiddenRoles {
		if strings.EqualFold(req.Role, forbiddenRole) {
			return nil, fmt.Errorf("registration with role '%s' is not allowed", req.Role)
		}
	}

	existingUser, err := s.Repo.GetUserByNameForClient(req.Name, int(req.ClientId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with given credentials already exists")
	}

	client, err := s.Repo.GetClientByID(int(req.ClientId), ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}
	if client == nil {
		return nil, errors.New("invalid client ID")
	}

	roleAllowed := false
	for _, role := range client.Roles {
		if strings.EqualFold(role, req.Role) {
			roleAllowed = true
			break
		}
	}
	if !roleAllowed {
		return nil, fmt.Errorf("role '%s' is not allowed for this client", req.Role)
	}

	hash, err := utils.GenerateHashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	log.Println("hash generated", hash)

	user := dto.User{
		ClientID: int(req.ClientId),
		Name:     req.Name,
		Password: hash,
		Role:     req.Role,
	}

	err = s.Repo.NewUser(&user)
	if err != nil {
		return nil, err
	}
	log.Println("user created: ", user)

	secret, err := utils.Decrypt(client.JwtSecret, config.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt client JWT secret: %v", err)
	}

	tokenString, err := utils.GenerateToken(user.ID, user.ClientID, user.Role, []byte(secret))
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Token: tokenString}, nil
}

func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.Repo.GetUserByNameForClient(req.Name, int(req.ClientId))
	if err != nil {
		return nil, errors.New("invalid name or password")
	}

	err = utils.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, errors.New("invalid name or password")
	}

	client, err := s.Repo.GetClientByID(user.ClientID, ctx)
	if client == nil {
		return nil, errors.New("invalid name or password")
	}

	secret, err := utils.Decrypt(client.JwtSecret, config.GetKey())

	tokenString, err := utils.GenerateToken(user.ID, user.ClientID, user.Role, []byte(secret))
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}

func (s *Service) RegisterAdmin(ctx context.Context, req *pb.RegisterAdminRequest, role string) (*pb.RegisterAdminResponse, error) {
	if !utils.HasAccess(role, dto.SupremeRole, dto.ClientSupremeRole) {
		return nil, errors.New("access denied, you are not supreme")
	}

	existingUser, err := s.Repo.GetUserByNameForClient(req.Name, int(req.ClientId))
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
		ClientID: int(req.ClientId),
		Name:     req.Name,
		Password: hash,
		Role:     req.Role,
	}

	err = s.Repo.NewUser(&user)
	if err != nil {
		return nil, err
	}
	log.Println("admin-user created: ", user)

	return &pb.RegisterAdminResponse{Id: int64(user.ID)}, nil
}

func (s *Service) RegisterClient(ctx context.Context, req *pb.RegisterClientRequest, role string) (*pb.Client, error) {
	utils.HasAccess(role, dto.SupremeRole, dto.ClientSupremeRole, dto.AdminRole)

	existingClient, err := s.Repo.GetClientByName(req.Name)
	if existingClient != nil {
		return nil, fmt.Errorf("user with given creds is existing: %w", err)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("registering client")
		} else {
			return nil, err
		}

	}
	code, err := utils.GenerateRandomString(6)
	if err != nil {
		return nil, err
	}

	encryptedCode, err := utils.Encrypt(code, config.GetKey())
	if err != nil {
		return nil, err
	}

	log.Println("hash generated", encryptedCode)

	client := dto.Client{
		Name:      req.Name,
		Roles:     req.Roles,
		JwtSecret: encryptedCode,
	}

	err = s.Repo.NewClient(&client)
	if err != nil {
		return nil, err
	}

	log.Println("client created: ", client)

	return &pb.Client{Id: int64(client.ID), Name: client.Name, Roles: client.Roles}, nil
}
