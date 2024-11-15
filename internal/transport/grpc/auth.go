package grpc

import (
	"context"
	"log"
	"service/internal/service"
	desc "service/pkg/grpc/auth_v1"
	pb "service/pkg/grpc/auth_v1"
)

type Server struct {
	desc.AuthServiceServer
	Service *service.Service
}

func NewServer(Service *service.Service) *Server {
	log.Printf("NewServer")
	return &Server{
		Service: Service,
	}
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return s.Service.Register(ctx, req)
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.Service.Login(ctx, req)
}
