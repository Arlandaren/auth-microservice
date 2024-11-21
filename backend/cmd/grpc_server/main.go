// server.go

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"service/internal/repository"
	"service/internal/service"
	"service/internal/shared/config"
	"service/internal/shared/logger"
	"service/internal/shared/storage/redis"

	log "github.com/sirupsen/logrus"

	prometheusModule "service/internal/shared/prometheus"
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
	logger.Init()
	ctx := context.Background()
	addresses := config.GetAddress()

	db, err := postgres.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB connected")

	rdb := redis.ConnectRedis(ctx)

	repo := repository.NewRepository(db, rdb)

	serv := service.NewService(repo)

	server := transport.NewServer(serv)

	prom := prometheusModule.NewPrometheus()
	prom.RegisterMetrics()

	wg := &sync.WaitGroup{}

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
		grpc.Creds(credentials),
		grpc.UnaryInterceptor(service.AuthInterceptor(server.Service.Repo, ctx)),
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

	c := allowCORS()

	handler := c.Handler(mux)

	wrappedHandler := prometheusModule.MetricsMiddleware(handler)

	httpsRouter := http.NewServeMux()
	httpsRouter.Handle("/", wrappedHandler)

	tlsConfig, err := utils.LoadServerTLS()
	if err != nil {
		return err
	}

	httpsServer := &http.Server{
		Addr:      addr.Http,
		Handler:   httpsRouter,
		TLSConfig: tlsConfig,
	}

	// Запуск HTTPS в отдельной горутине
	go func() {
		if err := httpsServer.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("HTTPS server exited with error: %v", err)
		}
	}()
	log.Infof("HTTPS server listening at %v\n", addr.Http)

	metricsRouter := http.NewServeMux()
	metricsRouter.Handle("/metrics", promhttp.Handler())

	metricsServerAddr := ":9000"

	metricsServer := &http.Server{
		Addr:    metricsServerAddr,
		Handler: metricsRouter,
	}

	// Запуск HTTP для метрик
	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("HTTP server for metrics exited with error: %v", err)
		}
	}()
	log.Infof("HTTP server for metrics listening at %v\n", metricsServerAddr)

	<-ctx.Done()

	log.Info("Shutting down servers...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpsServer.Shutdown(shutdownCtx); err != nil {
		log.Errorf("HTTPS server Shutdown failed: %v", err)
	}

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Errorf("HTTP server for metrics Shutdown failed: %v", err)
	}

	return nil
}
