// server.go

package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"service/internal/repository"
	"service/internal/service"
	"service/internal/shared/storage/postgres"
	pb "service/pkg/grpc/auth_v1"
	"sync"
	"syscall"

	transport "service/internal/transport/grpc"
)

const grpcAddress = ":8080"

func main() {
	db, err := postgres.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DB connected")
	repo := repository.NewRepository(db)

	serv := service.NewService(repo)

	server := transport.NewServer(serv)

	wg := &sync.WaitGroup{}
	ctx := context.Background()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := RunGrpcServer(ctx, server)
		if err != nil {
			log.Printf("failed to run grpc server: %v", err)
			return
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down servers...")

	wg.Wait()
	log.Println("Servers gracefully stopped.")
}

func RunGrpcServer(ctx context.Context, server *transport.Server) error {

	// Uncomment and configure the TLS credentials if needed
	// credentials, err := utils.LoadTLSCredentials()
	// if err != nil {
	//     return fmt.Errorf("failed to load TLS credentials: %w", err)
	// }

	grpcServer := grpc.NewServer(
		// grpc.Creds(credentials), // Uncomment if using TLS
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(service.AuthInterceptor),
	)
	reflection.Register(grpcServer)
	pb.RegisterAuthServiceServer(grpcServer, server)

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcAddress, err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Printf("gRPC server failed: %v", err)
		}
	}()

	log.Printf("gRPC server listening at %v\n", grpcAddress)

	<-ctx.Done()

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	return nil
}
