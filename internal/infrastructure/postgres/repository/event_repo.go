package postgres

import (
	"context"
	"database/sql"

	"shipment-service/internal/domain/shipment"
)

type StatusEventRepository struct {
	db *sql.DB
}

func NewStatusEventRepository(db *sql.DB) *StatusEventRepository {
	return &StatusEventRepository{db: db}
}

func (r *StatusEventRepository) Create(ctx context.Context, e *shipment.StatusEvent) error {
	query := `
		INSERT INTO shipment_events (id, shipment_id, status, note, occurred_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query, e.ID, e.ShipmentID, string(e.Status), e.Note, e.OccurredAt)
	return err
}

func (r *StatusEventRepository) GetAllByShipmentID(ctx context.Context, shipmentID string) ([]shipment.StatusEvent, error) {
	query := `
		SELECT id, shipment_id, status, note, occurred_at
		FROM shipment_events
		WHERE shipment_id = $1
		ORDER BY occurred_at ASC`

	rows, err := r.db.QueryContext(ctx, query, shipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []shipment.StatusEvent
	for rows.Next() {
		var e shipment.StatusEvent
		var status string
		if err := rows.Scan(&e.ID, &e.ShipmentID, &status, &e.Note, &e.OccurredAt); err != nil {
			return nil, err
		}
		e.Status = shipment.Status(status)
		events = append(events, e)
	}
	return events, rows.Err()
}
