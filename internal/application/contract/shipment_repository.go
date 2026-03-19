package contract

import (
	"context"

	"shipment-service/internal/domain/shipment"
)

type ShipmentRepository interface {
	Create(ctx context.Context, s *shipment.Shipment) error
	GetByID(ctx context.Context, id string) (*shipment.Shipment, error)
	GetByReferenceNumber(ctx context.Context, ref string) (*shipment.Shipment, error)
}
