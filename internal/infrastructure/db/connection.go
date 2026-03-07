package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	PINT_TIMEOUT          = 5 * time.Second
	DEFAULT_MAX_CONNS     = int32(25)
	DEFAULT_MIN_CONNS     = int32(2)
	DEFAULT_CONN_LIFETIME = 60 * time.Minute
)

func NewConnection(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {

	// parse db url and create pool configuration
	poolCnf, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse database config: %w", err)
	}

	// override default pool configuration
	poolCnf.MaxConns = DEFAULT_MAX_CONNS
	poolCnf.MinConns = DEFAULT_MIN_CONNS
	poolCnf.MaxConnLifetime = DEFAULT_CONN_LIFETIME

	// create pool with avobe configuration
	pool, err := pgxpool.NewWithConfig(ctx, poolCnf)
	if err != nil {
		return nil, fmt.Errorf("Unable to create database pool: %w", err)
	}

	// ping database
	if err := pingWithTimeout(ctx, pool); err != nil {
		pool.Close()
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	return pool, nil
}

func pingWithTimeout(ctx context.Context, pool *pgxpool.Pool) error {

	pingCtx, cancel := context.WithTimeout(ctx, PINT_TIMEOUT)
	defer cancel()

	err := pool.Ping(pingCtx)
	if err != nil {
		return fmt.Errorf("Unable to ping: %w", err)
	}

	return nil
}
