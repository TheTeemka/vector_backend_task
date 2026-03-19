package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"shipment-service/internal/application/contract/mocks"
	"shipment-service/internal/application/dto"
	"shipment-service/internal/domain"
	"shipment-service/internal/domain/shipment"
)

func newTestService(t *testing.T) (*ShipmentService, *mocks.ShipmentRepository, *mocks.StatusEventRepository, *mocks.IDGenerator) {
	shipmentRepo := mocks.NewShipmentRepository(t)
	eventRepo := mocks.NewStatusEventRepository(t)
	idGen := mocks.NewIDGenerator(t)
	txManager := mocks.NewTxManager(t)
	txManager.EXPECT().
		WithTx(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).Maybe()
	logger := zap.NewNop()
	svc := NewShipmentService(shipmentRepo, eventRepo, idGen, txManager, logger)
	return svc, shipmentRepo, eventRepo, idGen
}

func TestCreateShipment_Success(t *testing.T) {
	svc, shipmentRepo, eventRepo, idGen := newTestService(t)
	ctx := context.Background()

	idGen.EXPECT().NewID().Return("shipment-1").Once()
	idGen.EXPECT().NewID().Return("event-1").Once()
	shipmentRepo.EXPECT().GetByReferenceNumber(ctx, "REF-001").Return(nil, domain.ErrShipmentNotFound)
	shipmentRepo.EXPECT().Create(ctx, mock.AnythingOfType("*shipment.Shipment")).Return(nil)
	eventRepo.EXPECT().Create(ctx, mock.AnythingOfType("*shipment.StatusEvent")).Return(nil)

	input := dto.CreateShipmentInput{
		ReferenceNumber: "REF-001",
		Origin:          "Almaty",
		Destination:     "Astana",
		Driver:          shipment.DriverInfo{DriverID: "driver-1", UnitID: "unit-1"},
		ShipmentAmount:  100.0,
		DriverRevenue:   50.0,
	}

	sh, err := svc.CreateShipment(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, "shipment-1", sh.ID)
	assert.Equal(t, "REF-001", sh.ReferenceNumber)
	assert.Equal(t, shipment.StatusPending, sh.CurrentStatus)
}

func TestCreateShipment_DuplicateReference(t *testing.T) {
	svc, shipmentRepo, _, _ := newTestService(t)
	ctx := context.Background()

	existing := &shipment.Shipment{ID: "existing-1", ReferenceNumber: "REF-001"}
	shipmentRepo.EXPECT().GetByReferenceNumber(ctx, "REF-001").Return(existing, nil)

	input := dto.CreateShipmentInput{
		ReferenceNumber: "REF-001",
		Origin:          "Almaty",
		Destination:     "Astana",
		Driver:          shipment.DriverInfo{DriverID: "driver-1", UnitID: "unit-1"},
	}

	sh, err := svc.CreateShipment(ctx, input)

	assert.Nil(t, sh)
	assert.ErrorIs(t, err, domain.ErrDuplicateReference)
}

func TestGetShipment_Success(t *testing.T) {
	svc, shipmentRepo, _, _ := newTestService(t)
	ctx := context.Background()

	expected := &shipment.Shipment{ID: "shipment-1", ReferenceNumber: "REF-001"}
	shipmentRepo.EXPECT().GetByID(ctx, "shipment-1").Return(expected, nil)

	sh, err := svc.GetShipment(ctx, "shipment-1")

	require.NoError(t, err)
	assert.Equal(t, expected, sh)
}

func TestGetShipment_NotFound(t *testing.T) {
	svc, shipmentRepo, _, _ := newTestService(t)
	ctx := context.Background()

	shipmentRepo.EXPECT().GetByID(ctx, "shipment-1").Return(nil, domain.ErrShipmentNotFound)

	sh, err := svc.GetShipment(ctx, "shipment-1")

	assert.Nil(t, sh)
	assert.ErrorIs(t, err, domain.ErrShipmentNotFound)
}

func TestAddStatusEvent_Success(t *testing.T) {
	svc, shipmentRepo, eventRepo, idGen := newTestService(t)
	ctx := context.Background()

	existing := &shipment.Shipment{
		ID:            "shipment-1",
		CurrentStatus: shipment.StatusPending,
	}
	shipmentRepo.EXPECT().GetByID(ctx, "shipment-1").Return(existing, nil)
	shipmentRepo.EXPECT().Create(ctx, mock.AnythingOfType("*shipment.Shipment")).Return(nil)
	idGen.EXPECT().NewID().Return("event-1")
	eventRepo.EXPECT().Create(ctx, mock.AnythingOfType("*shipment.StatusEvent")).Return(nil)

	sh, err := svc.AddStatusEvent(ctx, "shipment-1", shipment.StatusPickedUp, "picked up")

	require.NoError(t, err)
	assert.Equal(t, shipment.StatusPickedUp, sh.CurrentStatus)
}

func TestAddStatusEvent_InvalidTransition(t *testing.T) {
	svc, shipmentRepo, _, idGen := newTestService(t)
	ctx := context.Background()

	existing := &shipment.Shipment{
		ID:            "shipment-1",
		CurrentStatus: shipment.StatusDelivered,
	}

	shipmentRepo.EXPECT().GetByID(ctx, "shipment-1").Return(existing, nil)
	idGen.EXPECT().NewID().Return("event-1")

	sh, err := svc.AddStatusEvent(ctx, "shipment-1", shipment.StatusInTransit, "try again")

	assert.Nil(t, sh)
	assert.ErrorIs(t, err, domain.ErrInvalidStatusTransition)
}

func TestGetShipmentHistory_Success(t *testing.T) {
	svc, shipmentRepo, eventRepo, _ := newTestService(t)
	ctx := context.Background()

	existing := &shipment.Shipment{ID: "shipment-1"}
	shipmentRepo.EXPECT().GetByID(ctx, "shipment-1").Return(existing, nil)

	events := []shipment.StatusEvent{
		{ID: "event-1", ShipmentID: "shipment-1", Status: shipment.StatusPending, Note: "created"},
		{ID: "event-2", ShipmentID: "shipment-1", Status: shipment.StatusInTransit, Note: "picked up"},
	}
	eventRepo.EXPECT().GetAllByShipmentID(ctx, "shipment-1").Return(events, nil)

	result, err := svc.GetShipmentHistory(ctx, "shipment-1")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "event-1", result[0].ID)
}
