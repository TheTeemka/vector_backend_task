package shipment

import (
	"fmt"
	"shipment-service/internal/domain"
	"time"
)

type StatusEvent struct {
	ID         string
	ShipmentID string
	Status     Status
	Note       string
	OccurredAt time.Time
}

func NewStatusEvent(id, shipmentID string, status Status, note string, occuredAt time.Time) (*StatusEvent, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if shipmentID == "" {
		return nil, fmt.Errorf("%w: shipment ID is required", domain.ErrInvalidInput)
	}

	if !IsValidStatus(status) {
		return nil, fmt.Errorf("%w: invalid status %s", domain.ErrInvalidInput, status)
	}

	return &StatusEvent{
		ID:         id,
		ShipmentID: shipmentID,
		Status:     status,
		Note:       note,
		OccurredAt: occuredAt,
	}, nil
}
