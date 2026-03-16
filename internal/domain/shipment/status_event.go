package shipment

import "time"

type StatusEvent struct {
	ID         string
	ShipmentID string
	Status     Status
	Note       string
	OccurredAt time.Time
}
