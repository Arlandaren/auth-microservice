// server.go

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"service/internal/repository"
	"service/internal/service"
	"service/internal/shared/config"
	"service/internal/shared/storage/dto"
	"service/internal/shared/storage/postgres"
	"service/internal/shared/utils"
	pb "service/pkg/grpc/auth_v1"
	"sync"
	"syscall"
	"time"

	transport "service/internal/transport/grpc"
)

func main() {
	//key, cert := tls.CreateAcCertificate()
	//tls.CreateServerCertificate(key, cert)

	addresses := config.GetAddress()

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
		err := RunGrpcServer(ctx, server, addresses)
		if err != nil {
			log.Printf("failed to run grpc server: %v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := startHttpServer(ctx, addresses)
		if err != nil {
			log.Printf("failed to run http server: %v", err)
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

func RunGrpcServer(ctx context.Context, server *transport.Server, addr *dto.Address) error {
	credentials, err := utils.LoadServerTLSCredentials()
	if err != nil {
		return fmt.Errorf("failed to load TLS credentials: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(credentials), // Uncomment if using TLS
		//grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(service.AuthInterceptor),
	)
	reflection.Register(grpcServer)
	pb.RegisterAuthServiceServer(grpcServer, server)

	lis, err := net.Listen("tcp", addr.Grpc)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr.Grpc, err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Printf("gRPC server failed: %v", err)
		}
	}()

	log.Printf("gRPC server listening at %v\n", addr.Grpc)

	<-ctx.Done()

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	return nil
}

func startHttpServer(ctx context.Context, addr *dto.Address) error {
	mux := runtime.NewServeMux()

	creds, err := utils.LoadClientTLSCredentials()
	if err != nil {
		return fmt.Errorf("failed to load client TLS credentials: %w", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	if err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, addr.Grpc, opts); err != nil {
		return fmt.Errorf("failed to register service handler: %w", err)
	}

	handler := allowCORS(mux)

	tlsConfig, err := utils.LoadServerTLS()
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:      addr.Http,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server exited with error: %v", err)
		}
	}()

	log.Printf("HTTP server listening at %v\n", addr.Http)

	<-ctx.Done()

	log.Println("Shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP server Shutdown failed: %w", err)
	}

	return nil
}
