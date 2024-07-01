package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juanfcgarcia/gostori/internal/repository/db"
)

type TransactionHandler struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewTransactionHandler(pool *pgxpool.Pool, queries *db.Queries) *TransactionHandler {
	return &TransactionHandler{
		pool:    pool,
		queries: queries,
	}
}

func (th *TransactionHandler) ExecuteInTransaction(ctx context.Context, fn func(pgx.Tx, *db.Queries) error) (err error) {
	tx, err := th.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				fmt.Println("Transaction rollback failed:", rollbackErr)
			}
		}
	}()

	queries := th.queries.WithTx(tx)

	err = fn(tx, queries)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		return
	}

	return
}
