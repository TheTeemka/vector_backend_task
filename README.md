# Shipment gRPC Microservice

A gRPC microservice for managing shipments and tracking status changes during transportation.

## API Reference

### gRPC Methods

| Method | Description |
|---|---|
| `CreateShipment` | Create a new shipment (starts as `pending`) |
| `GetShipment` | Retrieve shipment details by ID |
| `AddStatusEvent` | Transition shipment to a new status |
| `GetShipmentHistory` | Retrieve all status events for a shipment |


## How to Run

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [goose](https://github.com/pressly/goose) (for migrations)

### With Docker Compose

```bash
docker compose up --build
```

This starts both the PostgreSQL database and the gRPC server on port `50051`.

### Local Development

1. Start Docker:

```bash
docker compose up
```

2. Run database migrations:

```bash
goose up
```


The gRPC server listens on port `50051` by default (configurable via `GRPC_PORT` in `.env`).

### Environment Variables

| Variable | Description | Default |
|---|---|---|
| `POSTGRES_HOST` | Database host | `localhost` |
| `POSTGRES_PORT` | Database port | `5433` |
| `POSTGRES_USER` | Database user | `shipment` |
| `POSTGRES_PASSWORD` | Database password | `shipment123` |
| `POSTGRES_DB` | Database name | `shipment` |
| `POSTGRES_SSL_MODE` | SSL mode | `disable` |
| `GRPC_PORT` | gRPC server port | `50051` |
| `APP_ENV` | Application environment | `development` |

## How to Run Tests

```bash
go test ./... -v
```

Tests cover domain logic and service layer only — no database or gRPC server required.

### Regenerating Mocks

```bash
mockery --dir=internal/application/contract --all --output=internal/application/contract/mocks --outpkg=mocks --with-expecter
```

## Architecture Overview

The project follows **Clean Architecture** with three layers:

```
cmd/server/main.go          # Composition root
internal/
  domain/                   # Core business logic (no external dependencies`)
    shipment/               # Shipment entity, StatusEvent, DriverInfo, FSM
    errors.go               # Domain errors
  application/              # Use cases & orchestration
    contract/               # Interfaces
    service/                # ShipmentService (business orchestration)
    usecase/                # Individual use case definitions
    dto/                    # Data transfer objects
  infrastructure/           # Adapters (implements ports)
    grpc/                   # gRPC handler, mapper, interceptors, server
    postgres/               # PostgreSQL repository implementations
    logger/                 # Zap logger factory (env-based configuration)
    uuid/                   # UUID generator
  pkg/
    ctxlog/                 # Context-aware log field propagation
  config/                   # Environment-based configuration
init/
  proto/shipment/           # Protocol Buffers definition
  postgres_migrations/      # Goose SQL migrations
gen/proto/shipment/         # Generated protobuf Go code
```

## Design Decisions

### Status Lifecycle (FSM)

Shipment statuses follow a finite state machine:

```
pending --> picked_up --> in_transit --> delivered
  |            |              |
  v            v              v
cancelled  cancelled      cancelled
```

- `pending`: initial status, set on creation
- `picked_up`: driver has picked up the shipment
- `in_transit`: shipment is moving toward destination
- `delivered`: terminal state, shipment successfully delivered
- `cancelled`: terminal state, reachable from any non-terminal status

Invalid transitions are rejected with `ErrInvalidStatusTransition`. Self-transitions (e.g., `pending -> pending`) are also rejected.

### Status Events as History

Every status change creates a `StatusEvent` record with the new status, a note, and a timestamp. The shipment's `current_status` always reflects the latest applied event. An initial event with status `pending` is created automatically when a shipment is created.

## Assumptions

1. **Shipment IDs** are server-generated UUIDs, not client-provided.
2. **Reference numbers** must be unique across all shipments.
3. **Driver and unit details** are stored as embedded value objects, not separate entities.
4. **Cancelled and delivered** are terminal states — no further transitions allowed.
5. **Status events are append-only** — they cannot be modified or deleted.
6. **No authentication/authorization** is implemented (out of scope for this task).