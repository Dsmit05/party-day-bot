package api

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Dsmit05/party-day-bot/internal/logger"
	pb "github.com/Dsmit05/party-day-bot/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	cfg        ConfigI
	userServer pb.UserServer
	grpcServer *grpc.Server
	httpServer *http.Server
}

func NewServer(cfg ConfigI, bot BotI) *Server {
	userServer := NewUserServerAPI(bot)

	return &Server{cfg: cfg, userServer: userServer}
}

func (s *Server) StartGRPC() {
	listener, err := net.Listen("tcp", s.cfg.GetApiGRPCServerAddr())
	if err != nil {
		logger.Error("net.Listen() error", err)
		return
	}

	newGrpcServer := grpc.NewServer(grpc.ConnectionTimeout(s.cfg.GetApiGRPCServerTimeout()))
	s.grpcServer = newGrpcServer
	pb.RegisterUserServer(s.grpcServer, s.userServer)

	logger.Info("GRPC", "Start GRPC server")

	if err = s.grpcServer.Serve(listener); err != nil {
		logger.Error("grpcServer.Serve() error", err)
	}
}

func (s *Server) StartREST() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterUserHandlerFromEndpoint(ctx, mux, s.cfg.GetApiGRPCServerAddr(), opts); err != nil {
		logger.Error("RegisterUserHandlerFromEndpoint() error", err)

		return
	}

	newHTTPServer := &http.Server{
		Addr:    s.cfg.GetApiHTTPServerAddr(),
		Handler: mux,
	}
	s.httpServer = newHTTPServer

	logger.Info("HTTP", "Start HTTP server")

	if err := s.httpServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Info("ApiServer", "Server closed under request")
		} else {
			logger.Error("http.ListenAndServe() error", err)
		}
	}
}

func (s *Server) Stop() {
	stop := make(chan bool)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	go func() {
		s.httpServer.SetKeepAlivesEnabled(false)

		if err := s.httpServer.Shutdown(ctx); err != nil {
			logger.Error("Server Shutdown failed", err)

			return
		}

		s.grpcServer.GracefulStop()

		stop <- true
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
		logger.Error("Server context timeout", ctx.Err())
	case <-stop:
		logger.Info("Server", "Server closed under request")
	}

}
