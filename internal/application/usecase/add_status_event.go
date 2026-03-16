package usecase

import (
	"context"

	"shipment-service/internal/domain/shipment"
)

type statusEventAdder interface {
	AddStatusEvent(ctx context.Context, id string, status shipment.Status, note string) (*shipment.Shipment, error)
}

type AddStatusEventUseCase struct {
	service statusEventAdder
}

func NewAddStatusEventUseCase(service statusEventAdder) *AddStatusEventUseCase {
	return &AddStatusEventUseCase{service: service}
}

func (uc *AddStatusEventUseCase) Execute(ctx context.Context, id string, status shipment.Status, note string) (*shipment.Shipment, error) {
	return uc.service.AddStatusEvent(ctx, id, status, note)
}
