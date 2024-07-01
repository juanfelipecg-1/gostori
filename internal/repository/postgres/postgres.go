package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juanfcgarcia/gostori/internal/environment"
	"github.com/juanfcgarcia/gostori/internal/repository/db"
	"go.uber.org/zap"
)

type DB struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func NewDatabase(cfg *environment.Config, logger *zap.SugaredLogger) (*DB, error) {
	logger.Infow("connecting to database",
		"host", cfg.DBHost,
		"port", cfg.DBPort,
		"schema", cfg.DBSchema,
	)

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBSchema)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse environment: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	queries := db.New(pool)

	return &DB{Pool: pool, Queries: queries}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
