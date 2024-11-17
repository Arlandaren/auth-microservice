package oidc_auth

import (
	"context"
	pb "service/pkg/grpc/auth_v1"
)

type server struct {
	pb.UnimplementedAuthServiceServer // Встраивание базового типа
}

func (s *server) Login(ctx context.Context, req *pb.LoginOIDCRequest) (*pb.LoginOIDCResponse, error) {
	state := req.State // Используем состояние из запроса, но обычно лучше генерировать его
	// Здесь вы должны сформировать ваш OAuth2 URL
	url := "https://oauth2provider.com/auth?state=" + state // Здесь используйте вашу логику
	return &pb.LoginOIDCResponse{Url: url}, nil
}

func (s *server) Callback(ctx context.Context, req *pb.CallbackOIDCRequest) (*pb.CallbackOIDCResponse, error) {
	// Обрабатываем колбек и используем код и состояние для выполнения логики
	code := req.Code
	state := req.State
	// Здесь ваша логика для обработки кода, например, обмен кодом на токен

	return &pb.CallbackOIDCResponse{Message: "Callback processed successfully!"}, nil
}
