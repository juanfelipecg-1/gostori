package db_factory

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juanfcgarcia/gostori/internal/environment"
	"github.com/juanfcgarcia/gostori/internal/ports"
	"github.com/juanfcgarcia/gostori/internal/repository/db"
	"github.com/juanfcgarcia/gostori/internal/repository/postgres"
	"go.uber.org/zap"
)

const (
	PostgresDB = "postgresql"
)

type Repository struct {
	AccountRepository     ports.AccountRepository
	TransactionRepository ports.TransactionRepository
	TxHandler             *postgres.TransactionHandler
}

func NewRepository(ctx context.Context, cfg *environment.Config, logger *zap.SugaredLogger, existingPool *pgxpool.Pool) (*Repository, error) {
	switch cfg.DBType {
	case PostgresDB:
		return newPostgresRepository(ctx, cfg, logger, existingPool)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}
}

func newPostgresRepository(ctx context.Context, cfg *environment.Config, logger *zap.SugaredLogger, existingPool *pgxpool.Pool) (*Repository, error) {
	var dbQueries *db.Queries
	var pool *pgxpool.Pool

	if existingPool != nil {
		dbQueries = db.New(existingPool)
		pool = existingPool
	} else {
		dbInstance, err := postgres.NewDatabase(cfg, logger)
		if err != nil {
			logger.Errorf("failed to initialize PostgreSQL database: %v", err)
			return nil, fmt.Errorf("failed to initialize PostgreSQL database: %w", err)
		}
		dbQueries = dbInstance.Queries
		pool = dbInstance.Pool
	}

	accountRepo := postgres.NewAccountRepository(dbQueries)
	transactionRepo := postgres.NewTransactionRepository(dbQueries)
	txHandler := postgres.NewTransactionHandler(pool, dbQueries)

	return &Repository{
		AccountRepository:     accountRepo,
		TransactionRepository: transactionRepo,
		TxHandler:             txHandler,
	}, nil
}
