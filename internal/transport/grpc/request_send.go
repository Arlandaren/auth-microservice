package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"service/internal/shared/storage/dto"
	"service/internal/shared/utils"
	pb "service/pkg/grpc/auth_v1"
)

func SendClientRequest(ctx context.Context, target string, accessToken *dto.AccessToken) error {
	credentials, err := utils.LoadClientTLSCredentials()
	if err != nil {
		return fmt.Errorf("failed to load TLS credentials: %w", err)
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(credentials))
	if err != nil {
		return fmt.Errorf("не удалось установить соединение с сервером, err: %w", err)
	}
	defer func() { _ = conn.Close() }()

	client := pb.NewAuthServiceClient(conn)

	request := &pb.SendAccessTokenRequest{
		AccessToken: accessToken.Token,
		TokenType:   accessToken.TokenType,
		ExpiresAt:   accessToken.ExpiresAt.Unix(),
	}

	response, err := client.SendAccessToken(ctx, request)
	if err != nil {
		return fmt.Errorf("ошибка при вызове YourMethod: %w", err)
	}
	log.Printf("Ответ от сервера: %v", response)

	return nil
}
