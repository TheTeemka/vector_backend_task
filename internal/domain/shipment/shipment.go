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
	eventID string,
) (*Shipment, *StatusEvent, error) {
	if id == "" {
		return nil, nil, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if refNum == "" {
		return nil, nil, fmt.Errorf("%w: reference number is required", domain.ErrInvalidInput)
	}
	if origin == "" {
		return nil, nil, fmt.Errorf("%w: origin is required", domain.ErrInvalidInput)
	}
	if destination == "" {
		return nil, nil, fmt.Errorf("%w: destination is required", domain.ErrInvalidInput)
	}
	if driver.DriverID == "" {
		return nil, nil, fmt.Errorf("%w: driver ID is required", domain.ErrInvalidInput)
	}

	if eventID == "" {
		return nil, nil, fmt.Errorf("%w: event ID is required", domain.ErrInvalidInput)
	}

	now := time.Now()

	event, err := NewStatusEvent(eventID, id, StatusPending, "shipment created", now)
	if err != nil {
		return nil, nil, err
	}

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
	}, event, nil
}

func (s *Shipment) ApplyStatusEvent(eventID string, newStatus Status, note string) (*StatusEvent, error) {
	if !CanTransition(s.CurrentStatus, newStatus) {
		return nil, fmt.Errorf("%w: %s -> %s", domain.ErrInvalidStatusTransition, s.CurrentStatus, newStatus)
	}

	now := time.Now()
	s.CurrentStatus = newStatus
	s.UpdatedAt = now
	return NewStatusEvent(eventID, s.ID, newStatus, note, now)
}
