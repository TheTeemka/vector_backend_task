-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shipment_events (
    id VARCHAR(50) PRIMARY KEY,
    shipment_id VARCHAR(50) NOT NULL REFERENCES shipments(id),
    status VARCHAR(20) NOT NULL,
    note VARCHAR(500),
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_shipment_events_shipment_id ON shipment_events(shipment_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_shipment_events_shipment_id;
DROP TABLE IF EXISTS shipment_events;
-- +goose StatementEnd
