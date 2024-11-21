package config

import (
	"encoding/base64"
	"errors"
	log "github.com/sirupsen/logrus"
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

func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}

func GetRedis() *dto.RedisConfig {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	return &dto.RedisConfig{
		Addr:     redisAddress,
		Password: redisPassword,
	}
}

func GetKey() []byte {
	keyBase64 := os.Getenv("KEY")
	if keyBase64 == "" {
		log.Fatal("KEY environment variable is not set")
		return nil
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		log.Fatalf("invalid key: %v", err)
		return nil
	}

	if len(key) != 32 {
		log.Fatalf("invalid key size: %d bytes. Expected 32 bytes for AES-256.", len(key))
		return nil
	}

	return key
}
