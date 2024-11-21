package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"service/internal/repository"
	"service/internal/shared/utils"
	"strings"
)

func AuthInterceptor(repo *repository.Repository, ctx context.Context) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if info.FullMethod == "/auth_v1.AuthService/Login" || info.FullMethod == "/auth_v1.AuthService/Register" {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Printf("Missing meta")
			return nil, status.Error(codes.Unauthenticated, "missing meta")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			log.Printf("Missing token")
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		tokenStr := strings.Split(authHeader[0], " ")[1]

		claims, err := utils.ParseTokenWithoutVerification(tokenStr)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Invalid token ParseTokenWithoutVerification: %v", err))
		}

		secret, err := repo.GetClientJwtSecret(claims.ClientID, ctx)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Invalid token GetClientJwtSecret: %v", err))
		}

		fmt.Println("secret\n", secret)

		claims, err = utils.ValidateToken(tokenStr, []byte(secret))
		if err != nil {
			log.Printf("Invalid token: %v", err)
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Invalid token ValidateToken: %v", err))
		}

		newCtx := context.WithValue(ctx, "userID", claims.UserID)
		newCtx = context.WithValue(newCtx, "role", claims.Role)

		return handler(newCtx, req)
	}
}
