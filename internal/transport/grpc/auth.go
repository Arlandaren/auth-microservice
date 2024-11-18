package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
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

func (s *Server) RegisterAdmin(ctx context.Context, req *pb.RegisterAdminRequest) (*pb.RegisterAdminResponse, error) {
	roleValue := ctx.Value("role")
	if roleValue == nil {
		return nil, status.Error(codes.Unauthenticated, "Undefined role")
	}
	role, ok := roleValue.(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Invalid role type")
	}

	return s.Service.RegisterAdmin(ctx, req, role)
}

func (s *Server) LoginOIDC(ctx context.Context, req *pb.LoginOIDCRequest) (*pb.LoginOIDCResponse, error) {
	return s.Service.LoginOIDC(ctx, req)
}

func (s *Server) Callback(ctx context.Context, req *pb.CallbackOIDCRequest) (*pb.CallbackOIDCResponse, error) {
	return s.Service.Callback(ctx, req)
}
