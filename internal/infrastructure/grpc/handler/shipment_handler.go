package handler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "shipment-service/gen/proto/shipment"
	"shipment-service/internal/application/dto"
	"shipment-service/internal/application/usecase"
	"shipment-service/internal/domain"
	"shipment-service/internal/infrastructure/grpc/mapper"
)

type ShipmentHandler struct {
	pb.UnimplementedShipmentServiceServer

	createShipment     *usecase.CreateShipmentUseCase
	getShipment        *usecase.GetShipmentUseCase
	addStatusEvent     *usecase.AddStatusEventUseCase
	getShipmentHistory *usecase.GetShipmentHistoryUseCase
	log                *zap.Logger
}

func NewShipmentHandler(
	createShipment *usecase.CreateShipmentUseCase,
	getShipment *usecase.GetShipmentUseCase,
	addStatusEvent *usecase.AddStatusEventUseCase,
	getShipmentHistory *usecase.GetShipmentHistoryUseCase,
	log *zap.Logger,
) *ShipmentHandler {
	return &ShipmentHandler{
		createShipment:     createShipment,
		getShipment:        getShipment,
		addStatusEvent:     addStatusEvent,
		getShipmentHistory: getShipmentHistory,
		log:                log,
	}
}

func (h *ShipmentHandler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	input := dto.CreateShipmentInput{
		ReferenceNumber: req.ReferenceNumber,
		Origin:          req.Origin,
		Destination:     req.Destination,
		Driver:          mapper.ProtoToDriverInfo(req.Driver),
		ShipmentAmount:  req.ShipmentAmount,
		DriverRevenue:   req.DriverRevenue,
	}

	s, err := h.createShipment.Execute(ctx, input)
	if err != nil {
		return nil, h.toGRPCError(err)
	}

	return &pb.CreateShipmentResponse{Shipment: mapper.ShipmentToProto(s)}, nil
}

func (h *ShipmentHandler) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error) {
	s, err := h.getShipment.Execute(ctx, req.Id)
	if err != nil {
		return nil, h.toGRPCError(err)
	}

	return &pb.GetShipmentResponse{Shipment: mapper.ShipmentToProto(s)}, nil
}

func (h *ShipmentHandler) AddStatusEvent(ctx context.Context, req *pb.AddStatusEventRequest) (*pb.AddStatusEventResponse, error) {
	domainStatus, err := mapper.ProtoToStatus(req.Status)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	s, err := h.addStatusEvent.Execute(ctx, req.Id, domainStatus, req.Note)
	if err != nil {
		return nil, h.toGRPCError(err)
	}

	return &pb.AddStatusEventResponse{Shipment: mapper.ShipmentToProto(s)}, nil
}

func (h *ShipmentHandler) GetShipmentHistory(ctx context.Context, req *pb.GetShipmentHistoryRequest) (*pb.GetShipmentHistoryResponse, error) {
	events, err := h.getShipmentHistory.Execute(ctx, req.Id)
	if err != nil {
		return nil, h.toGRPCError(err)
	}

	pbEvents := make([]*pb.StatusEvent, 0, len(events))
	for _, e := range events {
		pbEvents = append(pbEvents, mapper.StatusEventToProto(e))
	}

	return &pb.GetShipmentHistoryResponse{Events: pbEvents}, nil
}

func (h *ShipmentHandler) toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrShipmentNotFound):
		return status.Errorf(codes.NotFound, "%v", err)
	case errors.Is(err, domain.ErrInvalidStatusTransition):
		return status.Errorf(codes.FailedPrecondition, "%v", err)
	case errors.Is(err, domain.ErrDuplicateReference):
		return status.Errorf(codes.AlreadyExists, "%v", err)
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Errorf(codes.InvalidArgument, "%v", err)
	default:
		return status.Errorf(codes.Internal, "internal server error")
	}
}
