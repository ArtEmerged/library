package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ArtEmerged/library/client/db"
)

type pgClient struct {
	masterDBC db.DB
}

// New creates a new instance of db.Client
func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connection to postgres db: %w", err)
	}

	return &pgClient{masterDBC: &pg{dbc: dbc}}, nil
}

// DB returns db.DB
func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

// Close closes db.DB
func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
