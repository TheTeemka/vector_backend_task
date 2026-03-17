-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shipments (
    id VARCHAR(50) PRIMARY KEY,
    reference_number VARCHAR(50) UNIQUE NOT NULL,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    current_status VARCHAR(20) NOT NULL,
    driver_id VARCHAR(50) NOT NULL,
    unit_id VARCHAR(50) NOT NULL,
    shipment_amount DECIMAL(10,2) NOT NULL,
    driver_revenue DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_shipments_reference_number ON shipments(reference_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_shipments_reference_number;
DROP TABLE IF EXISTS shipments;
-- +goose StatementEnd
