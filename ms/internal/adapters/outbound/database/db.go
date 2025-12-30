package database

import (
	"context"
	"fmt"

	database "github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

type Store struct {
	*DB
	*database.Queries
}

func NewStore(dbConn *DB, queires *database.Queries) *Store {
	return &Store{
		dbConn,
		queires,
	}
}

func NewDB(ctx context.Context, dsn string) (*DB, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
