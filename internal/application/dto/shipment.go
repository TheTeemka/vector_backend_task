package dto

import "shipment-service/internal/domain/shipment"

type CreateShipmentInput struct {
	ReferenceNumber string              `validate:"required"`
	Origin          string              `validate:"required"`
	Destination     string              `validate:"required"`
	Driver          shipment.DriverInfo `validate:"required"`
	ShipmentAmount  float64             `validate:"gte=0"`
	DriverRevenue   float64             `validate:"gte=0"`
}

func (i CreateShipmentInput) ToEntity(id string) (*shipment.Shipment, error) {
	return shipment.NewShipment(
		id,
		i.ReferenceNumber,
		i.Origin,
		i.Destination,
		i.Driver,
		i.ShipmentAmount,
		i.DriverRevenue,
	)
}
