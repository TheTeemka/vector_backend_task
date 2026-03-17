package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"shipment-service/internal/application/service"
	"shipment-service/internal/application/usecase"
	"shipment-service/internal/config"
	"shipment-service/internal/infrastructure/grpc/handler"
	"shipment-service/internal/infrastructure/grpc/server"
	zaplogger "shipment-service/internal/infrastructure/logger"
	"shipment-service/internal/infrastructure/postgres"
	postgresrepo "shipment-service/internal/infrastructure/postgres/repository"
	"shipment-service/internal/infrastructure/uuid"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger, err := zaplogger.New(zaplogger.Environment(cfg.AppEnv))
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer logger.Sync()
	logger.Info("starting server", zap.String("env", cfg.AppEnv), zap.String("port", cfg.GRPCPort))

	db, err := postgres.NewDB(cfg.Database.DSN())
	if err != nil {
		logger.Fatal("connect to postgres", zap.Error(err))
	}
	defer db.Close()

	shipmentRepo := postgresrepo.NewShipmentRepository(db)
	eventRepo := postgresrepo.NewStatusEventRepository(db)
	idGen := uuid.NewGenerator()

	shipmentSVC := service.NewShipmentService(shipmentRepo, eventRepo, idGen, logger)

	h := handler.NewShipmentHandler(
		usecase.NewCreateShipmentUseCase(shipmentSVC),
		usecase.NewGetShipmentUseCase(shipmentSVC),
		usecase.NewAddStatusEventUseCase(shipmentSVC),
		usecase.NewGetShipmentHistoryUseCase(shipmentSVC),
		logger,
	)

	srv := server.New(cfg.GRPCPort, h, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(); err != nil {
			logger.Fatal("grpc server error", zap.Error(err))
		}
	}()

	<-quit
	srv.GracefulStop()
}
