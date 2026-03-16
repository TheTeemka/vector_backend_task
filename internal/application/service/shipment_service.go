package service

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"shipment-service/internal/application/contract"
	"shipment-service/internal/application/dto"
	"shipment-service/internal/domain"
	"shipment-service/internal/domain/shipment"
)

type ShipmentService struct {
	shipmentRepo contract.ShipmentRepository
	eventRepo    contract.StatusEventRepository
	idGen        contract.IDGenerator
	logger       *zap.Logger
	validate     *validator.Validate
}

func NewShipmentService(
	shipmentRepo contract.ShipmentRepository,
	eventRepo contract.StatusEventRepository,
	idGen contract.IDGenerator,
	logger *zap.Logger,
) *ShipmentService {
	return &ShipmentService{
		shipmentRepo: shipmentRepo,
		eventRepo:    eventRepo,
		idGen:        idGen,
		logger:       logger,
		validate:     validator.New(),
	}
}

func (s *ShipmentService) CreateShipment(ctx context.Context, input dto.CreateShipmentInput) (*shipment.Shipment, error) {
	if err := s.validate.Struct(input); err != nil {
		return nil, err
	}

	existing, err := s.shipmentRepo.GetByReferenceNumber(ctx, input.ReferenceNumber)
	if err != nil && !errors.Is(err, domain.ErrShipmentNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrDuplicateReference
	}

	sh, err := input.ToEntity(s.idGen.NewID())
	if err != nil {
		return nil, err
	}

	if err := s.shipmentRepo.Save(ctx, sh); err != nil {
		return nil, err
	}

	initialEvent := &shipment.StatusEvent{
		ID:         s.idGen.NewID(),
		ShipmentID: sh.ID,
		Status:     shipment.StatusPending,
		Note:       "shipment created",
	}
	if err := s.eventRepo.Create(ctx, initialEvent); err != nil {
		return nil, err
	}

	return sh, nil
}

func (s *ShipmentService) GetShipment(ctx context.Context, id string) (*shipment.Shipment, error) {
	return s.shipmentRepo.GetByID(ctx, id)
}

func (s *ShipmentService) AddStatusEvent(ctx context.Context, id string, status shipment.Status, note string) (*shipment.Shipment, error) {
	sh, err := s.shipmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	event, err := sh.ApplyStatusEvent(s.idGen.NewID(), status, note)
	if err != nil {
		return nil, err
	}

	if err := s.shipmentRepo.Save(ctx, sh); err != nil {
		return nil, err
	}
	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return sh, nil
}

func (s *ShipmentService) GetShipmentHistory(ctx context.Context, id string) ([]shipment.StatusEvent, error) {
	if _, err := s.shipmentRepo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	return s.eventRepo.GetAllByShipmentID(ctx, id)
}
