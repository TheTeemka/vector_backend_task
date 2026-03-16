package contract

import (
	"context"

	"shipment-service/internal/domain/shipment"
)

type StatusEventRepository interface {
	Create(ctx context.Context, e *shipment.StatusEvent) error
	GetAllByShipmentID(ctx context.Context, shipmentID string) ([]shipment.StatusEvent, error)
}
