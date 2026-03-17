package shipment

import (
	"fmt"
	"time"

	"shipment-service/internal/domain"
)

type Shipment struct {
	ID              string
	ReferenceNumber string
	Origin          string
	Destination     string
	CurrentStatus   Status
	Driver          DriverInfo
	ShipmentAmount  float64
	DriverRevenue   float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewShipment(
	id, refNum, origin, destination string,
	driver DriverInfo,
	amount, revenue float64,
) (*Shipment, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if refNum == "" {
		return nil, fmt.Errorf("%w: reference number is required", domain.ErrInvalidInput)
	}
	if origin == "" {
		return nil, fmt.Errorf("%w: origin is required", domain.ErrInvalidInput)
	}
	if destination == "" {
		return nil, fmt.Errorf("%w: destination is required", domain.ErrInvalidInput)
	}
	if driver.DriverID == "" {
		return nil, fmt.Errorf("%w: driver ID is required", domain.ErrInvalidInput)
	}

	now := time.Now()
	return &Shipment{
		ID:              id,
		ReferenceNumber: refNum,
		Origin:          origin,
		Destination:     destination,
		CurrentStatus:   StatusPending,
		Driver:          driver,
		ShipmentAmount:  amount,
		DriverRevenue:   revenue,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (s *Shipment) ApplyStatusEvent(eventID string, newStatus Status, note string) (*StatusEvent, error) {
	if !CanTransition(s.CurrentStatus, newStatus) {
		return nil, fmt.Errorf("%w: %s -> %s", domain.ErrInvalidStatusTransition, s.CurrentStatus, newStatus)
	}

	now := time.Now()
	s.CurrentStatus = newStatus
	s.UpdatedAt = now
	return &StatusEvent{
		ID:         eventID,
		ShipmentID: s.ID,
		Status:     newStatus,
		Note:       note,
		OccurredAt: now,
	}, nil
}
