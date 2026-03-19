package dto

import "shipment-service/internal/domain/shipment"

type CreateShipmentInput struct {
	ReferenceNumber string              `validate:"required"`
	Origin          string              `validate:"required"`
	Destination     string              `validate:"required"`
	Driver          shipment.DriverInfo `validate:"required"`
	ShipmentAmount  float64
	DriverRevenue   float64
}

func (i CreateShipmentInput) ToEntity(id string, eventID string) (*shipment.Shipment, *shipment.StatusEvent, error) {
	return shipment.NewShipment(
		id,
		i.ReferenceNumber,
		i.Origin,
		i.Destination,
		i.Driver,
		i.ShipmentAmount,
		i.DriverRevenue,
		eventID,
	)
}
