package server

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "shipment-service/gen/proto/shipment"
	"shipment-service/internal/infrastructure/grpc/handler"
	"shipment-service/internal/infrastructure/grpc/interceptor"
)

type Server struct {
	grpc   *grpc.Server
	port   string
	logger *zap.Logger
}

func New(port string, h *handler.ShipmentHandler, log *zap.Logger) *Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryRecovery(log),
			interceptor.UnaryLogger(log),
		),
	)
	pb.RegisterShipmentServiceServer(srv, h)

	return &Server{
		grpc:   srv,
		port:   port,
		logger: log,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	return s.grpc.Serve(lis)
}

func (s *Server) GracefulStop() {
	s.grpc.GracefulStop()
}
