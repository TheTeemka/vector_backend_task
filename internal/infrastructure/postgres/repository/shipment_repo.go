package postgres

import (
	"context"
	"database/sql"
	"errors"

	"shipment-service/internal/domain"
	"shipment-service/internal/domain/shipment"
	"shipment-service/internal/infrastructure/postgres"
)

type ShipmentRepository struct {
	db *sql.DB
}

func NewShipmentRepository(db *sql.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) Save(ctx context.Context, s *shipment.Shipment) error {
	query := `
		INSERT INTO shipments (id, reference_number, origin, destination, current_status,
		                       driver_id, unit_id, shipment_amount, driver_revenue, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			current_status = EXCLUDED.current_status,
			updated_at     = EXCLUDED.updated_at`

	exec := postgres.ExtractExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, query,
		s.ID, s.ReferenceNumber, s.Origin, s.Destination, string(s.CurrentStatus),
		s.Driver.DriverID, s.Driver.UnitID,
		s.ShipmentAmount, s.DriverRevenue,
		s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *ShipmentRepository) GetByID(ctx context.Context, id string) (*shipment.Shipment, error) {
	query := `
		SELECT id, reference_number, origin, destination, current_status,
		       driver_id, unit_id, shipment_amount, driver_revenue, created_at, updated_at
		FROM shipments WHERE id = $1`

	exec := postgres.ExtractExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, id)
	return scanShipment(row)
}

func (r *ShipmentRepository) GetByReferenceNumber(ctx context.Context, ref string) (*shipment.Shipment, error) {
	query := `
		SELECT id, reference_number, origin, destination, current_status,
		       driver_id, unit_id, shipment_amount, driver_revenue, created_at, updated_at
		FROM shipments WHERE reference_number = $1`

	exec := postgres.ExtractExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, ref)
	return scanShipment(row)
}

func scanShipment(row *sql.Row) (*shipment.Shipment, error) {
	var s shipment.Shipment
	var status string
	err := row.Scan(
		&s.ID, &s.ReferenceNumber, &s.Origin, &s.Destination, &status,
		&s.Driver.DriverID, &s.Driver.UnitID,
		&s.ShipmentAmount, &s.DriverRevenue,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrShipmentNotFound
	}
	if err != nil {
		return nil, err
	}
	s.CurrentStatus = shipment.Status(status)
	return &s, nil
}
