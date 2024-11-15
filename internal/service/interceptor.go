package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"service/internal/shared/utils"
	"strings"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/auth_v1.AuthService/Login" || info.FullMethod == "/auth_v1.AuthService/Register" {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing meta")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	tokenStr := strings.Split(authHeader[0], " ")[1]

	claims, err := utils.ValidateToken(tokenStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Invalid token: %v", err))
	}

	newCtx := context.WithValue(ctx, "userID", claims.UserID)
	newCtx = context.WithValue(newCtx, "role", claims.Role)

	return handler(newCtx, req)
}
