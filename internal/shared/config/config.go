package config

import (
	"errors"
	"fmt"
	"os"
	"service/internal/shared/storage/dto"
)

func GetPostgres() (*dto.PostgresConfig, error) {
	pgConn := os.Getenv("PG_STRING")
	if pgConn == "" {
		return nil, errors.New("not found PG_STRING")
	}
	return &dto.PostgresConfig{
		ConnStr: pgConn,
	}, nil
}

func GetJwt() string {
	jwt := os.Getenv("jwt_key")
	fmt.Println(jwt)
	if jwt == "" {
		return "secret"
	}
	return jwt
}

func GetAddress() *dto.Address {
	httpAddress := os.Getenv("HTTP_ADDRESS")
	grpcAddress := os.Getenv("GRPC_ADDRESS")

	if httpAddress == "" {
		httpAddress = ":8086"
	}
	if grpcAddress == "" {
		grpcAddress = ":50051"
	}

	return &dto.Address{
		Http: httpAddress,
		Grpc: grpcAddress,
	}
}

func GetAccessData() *dto.AccessData {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectUrl := os.Getenv("REDIRECT_URL")

	return &dto.AccessData{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
	}
}
