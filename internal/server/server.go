package server

import (
	"context"
	"fmt"
	"github.com/codejago/polypully-openai-checker/internal/cache"
	"github.com/codejago/polypully-openai-checker/internal/metrics"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"time"

	pb "github.com/codejago/polypully-openai-checker/api/software"
	"github.com/codejago/polypully-openai-checker/internal/config"
	"github.com/codejago/polypully-openai-checker/internal/openai"
	"github.com/codejago/polypully-openai-checker/internal/security"
	"github.com/codejago/polypully-openai-checker/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{config: config}
}

func (s *Server) StartServer(systemMessage *openai.Message, oai *openai.OpenAI) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Server.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	// new gRPC server
	grpcServer, err := s.createGRPCServer()
	if err != nil {
		return fmt.Errorf("failed to create gRPC server: %w", err)
	}
	// cache will never be nil but may be disabled depending on the config
	cachePtr := cache.NewCache(s.config)
	// register the service with configuration
	pb.RegisterSoftwareInfoServiceServer(grpcServer, &service.Service{
		SystemMessage: *systemMessage,
		Oai:           oai,
		Config:        *s.config,
		Cache:         cachePtr,
	})
	reflection.Register(grpcServer)

	slog.Info("server listening", "address", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *Server) createGRPCServer() (*grpc.Server, error) {
	if s.config.Tls.Enabled {
		tls := security.NewTls()
		tlsConfig, err := tls.SetupTLS(s.config.Tls)
		if err != nil {
			return nil, fmt.Errorf("failed to setup TLS: %w", err)
		}
		creds := credentials.NewTLS(tlsConfig)
		return grpc.NewServer(grpc.Creds(creds), grpc.ChainUnaryInterceptor(metricsUnaryInterceptor)), nil
	}
	return grpc.NewServer(grpc.ChainUnaryInterceptor(metricsUnaryInterceptor)), nil
}

func metricsUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	code := status.Code(err)
	metrics.IncRequestCounter(info.FullMethod, code.String())
	metrics.SetRequestDuration(info.FullMethod, code.String(), time.Since(start).Seconds())
	return resp, err
}
