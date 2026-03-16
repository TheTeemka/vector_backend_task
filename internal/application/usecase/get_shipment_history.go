package usecase

import (
	"context"

	"shipment-service/internal/domain/shipment"
)

type shipmentHistoryGetter interface {
	GetShipmentHistory(ctx context.Context, id string) ([]shipment.StatusEvent, error)
}

type GetShipmentHistoryUseCase struct {
	service shipmentHistoryGetter
}

func NewGetShipmentHistoryUseCase(service shipmentHistoryGetter) *GetShipmentHistoryUseCase {
	return &GetShipmentHistoryUseCase{service: service}
}

func (uc *GetShipmentHistoryUseCase) Execute(ctx context.Context, id string) ([]shipment.StatusEvent, error) {
	return uc.service.GetShipmentHistory(ctx, id)
}
