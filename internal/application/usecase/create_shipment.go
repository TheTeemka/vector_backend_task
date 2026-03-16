package usecase

import (
	"context"

	"shipment-service/internal/application/dto"
	"shipment-service/internal/domain/shipment"
)

type shipmentCreator interface {
	CreateShipment(ctx context.Context, input dto.CreateShipmentInput) (*shipment.Shipment, error)
}

type CreateShipmentUseCase struct {
	service shipmentCreator
}

func NewCreateShipmentUseCase(service shipmentCreator) *CreateShipmentUseCase {
	return &CreateShipmentUseCase{service: service}
}

func (uc *CreateShipmentUseCase) Execute(ctx context.Context, input dto.CreateShipmentInput) (*shipment.Shipment, error) {
	return uc.service.CreateShipment(ctx, input)
}
