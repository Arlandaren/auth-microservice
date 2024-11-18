package test

import (
	"google.golang.org/grpc"
	"log"
	"service/internal/shared/utils"
)

func main() {
	creds, err := utils.LoadClientTLSCredentials()
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}

	connect, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer func() { _ = connect.Close() }()

}
