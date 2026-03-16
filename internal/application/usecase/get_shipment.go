package usecase

import (
	"context"

	"shipment-service/internal/domain/shipment"
)

type shipmentGetter interface {
	GetShipment(ctx context.Context, id string) (*shipment.Shipment, error)
}

type GetShipmentUseCase struct {
	service shipmentGetter
}

func NewGetShipmentUseCase(service shipmentGetter) *GetShipmentUseCase {
	return &GetShipmentUseCase{service: service}
}

func (uc *GetShipmentUseCase) Execute(ctx context.Context, id string) (*shipment.Shipment, error) {
	return uc.service.GetShipment(ctx, id)
}
