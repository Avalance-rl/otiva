package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool struct {
	Pool *pgxpool.Pool
}

func NewPGX(
	host, port, sslMode, user, password, name string,
	maxConns int32,
) (*PgxPool, error) {
	config, err := pgxpool.ParseConfig(buildPostgresDns(host, port, sslMode, user, password, name))
	if err != nil {
		return nil, err
	}
	config.MaxConns = maxConns

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PgxPool{pool}, nil
}

func Close(ctx context.Context, pool *pgxpool.Pool) error {
	done := make(chan struct{}, 1)
	go func() {
		pool.Close()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func buildPostgresDns(host, port, sslMode, user, password, name string) string {
	return fmt.Sprintf(
		"host=%s port=%s sslmode=%s user=%s password=%s dbname=%s",
		host, port, sslMode, user, password, name,
	)
}
