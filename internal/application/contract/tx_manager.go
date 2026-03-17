package contract

import "context"

// TxManager runs fn inside a database transaction.
// If fn returns nil the transaction is committed; otherwise it is rolled back.
type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
