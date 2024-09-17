package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Handler is a function that will be executed in transaction
type Handler func(context.Context) error

// Query is a query
type Query struct {
	Name     string
	QueryRaw string
}

// Client is an interface for db client
type Client interface {
	DB() DB
	Close() error
}

// DB is an interface for db
type DB interface {
	SQLExec
	Transaction
	Pinger
	Close()
}

// SQLExec executes query
type SQLExec interface {
	NamedExecer
	QueryExecer
}

// NamedExecer executes query
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest any, query Query, args ...any) error
	ScanAllContext(ctx context.Context, dest any, query Query, args ...any) error
}

// QueryExecer executes query
type QueryExecer interface {
	ExecContext(ctx context.Context, query Query, args ...any) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, query Query, args ...any) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, query Query, args ...any) pgx.Row
}

// Pinger checks connection
type Pinger interface {
	Ping(ctx context.Context) error
}

// Transaction is a transaction interface
type Transaction interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// TxManager is a transaction manager
type TxManager interface {
	ReadCommitted(ctx context.Context, fn Handler) error
}
