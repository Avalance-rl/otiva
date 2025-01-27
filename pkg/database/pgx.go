package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Pgx struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPGX(
	host, port, sslMode, user, password, name string,
	maxConns int32,
	logger *zap.Logger,
) (*Pgx, error) {
	config, err := pgxpool.ParseConfig(buildPostgresDns(host, port, sslMode, user, password, name))
	if err != nil {
		return nil, err
	}
	config.MaxConns = maxConns

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Pgx{db, logger}, nil
}

func (p *Pgx) Close(ctx context.Context) error {
	done := make(chan struct{}, 1)
	go func() {
		p.pool.Close()
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
