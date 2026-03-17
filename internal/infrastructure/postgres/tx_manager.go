package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

// DBTX is the interface satisfied by both *sql.DB and *sql.Tx.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type ctxKeyTx struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

// ExtractExecutor returns the *sql.Tx stored in ctx, or falls back to db.
func ExtractExecutor(ctx context.Context, db *sql.DB) DBTX {
	if tx, ok := ctx.Value(ctxKeyTx{}).(*sql.Tx); ok {
		return tx
	}
	return db
}

// TxManager implements contract.TxManager using database/sql transactions.
type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = fn(injectTx(ctx, tx))
	if err != nil {
		return err
	}
	return tx.Commit()
}
