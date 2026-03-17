package shipment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanTransition_ValidTransitions(t *testing.T) {
	tests := []struct {
		name string
		from Status
		to   Status
	}{
		{"pending to in_transit", StatusPending, StatusInTransit},
		{"pending to cancelled", StatusPending, StatusCancelled},
		{"in_transit to delivered", StatusInTransit, StatusDelivered},
		{"in_transit to cancelled", StatusInTransit, StatusCancelled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, CanTransition(tt.from, tt.to))
		})
	}
}

func TestCanTransition_InvalidTransitions(t *testing.T) {
	tests := []struct {
		name string
		from Status
		to   Status
	}{
		{"pending to delivered", StatusPending, StatusDelivered},
		{"in_transit to pending", StatusInTransit, StatusPending},
		{"delivered to pending", StatusDelivered, StatusPending},
		{"delivered to in_transit", StatusDelivered, StatusInTransit},
		{"delivered to cancelled", StatusDelivered, StatusCancelled},
		{"cancelled to pending", StatusCancelled, StatusPending},
		{"cancelled to in_transit", StatusCancelled, StatusInTransit},
		{"cancelled to delivered", StatusCancelled, StatusDelivered},
		{"same status pending", StatusPending, StatusPending},
		{"same status in_transit", StatusInTransit, StatusInTransit},
		{"same status delivered", StatusDelivered, StatusDelivered},
		{"same status cancelled", StatusCancelled, StatusCancelled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.False(t, CanTransition(tt.from, tt.to))
		})
	}
}

func TestCanTransition_UnknownStatus(t *testing.T) {
	assert.False(t, CanTransition(Status("unknown"), StatusPending))
	assert.False(t, CanTransition(StatusPending, Status("unknown")))
}
