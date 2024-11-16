package service

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	"service/internal/repository"
	"service/internal/shared/storage/dto"
	"service/internal/shared/utils"
	pb "service/pkg/grpc/auth_v1"
	"strings"
)

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
