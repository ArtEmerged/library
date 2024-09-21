package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ArtEmerged/library/client/db"
	"github.com/ArtEmerged/library/client/db/prettier"
)

type key string

const (
	// TxKey is a context key for pgx.Tx
	TxKey key = "tx"
)

type pg struct {
	dbc *pgxpool.Pool
}

// NewDB creates a new instance of db.DB
func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{dbc: dbc}
}

// ScanOneContext executes query and scans result into dest
func (p *pg) ScanOneContext(ctx context.Context, dest any, query db.Query, args ...any) error {
	logQuery(ctx, query, args...)

	row, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext executes query and scans result into dest
func (p *pg) ScanAllContext(ctx context.Context, dest any, query db.Query, args ...any) error {
	logQuery(ctx, query, args...)

	rows, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

// ExecContext executes query
func (p *pg) ExecContext(ctx context.Context, query db.Query, args ...any) (pgconn.CommandTag, error) {
	logQuery(ctx, query, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, query.QueryRaw, args...)
	}

	return p.dbc.Exec(ctx, query.QueryRaw, args...)
}

// QueryContext executes query
func (p *pg) QueryContext(ctx context.Context, query db.Query, args ...any) (pgx.Rows, error) {
	logQuery(ctx, query, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, query.QueryRaw, args...)
	}

	return p.dbc.Query(ctx, query.QueryRaw, args...)
}

// QueryRowContext executes query
func (p *pg) QueryRowContext(ctx context.Context, query db.Query, args ...any) pgx.Row {
	logQuery(ctx, query, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query.QueryRaw, args...)
	}

	return p.dbc.QueryRow(ctx, query.QueryRaw, args...)
}

// BeginTx starts transaction
func (p *pg) BeginTx(ctx context.Context, txOption pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOption)
}

// Ping checks connection
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close closes connection
func (p *pg) Close() {
	p.dbc.Close()
}

// MakeContextTx creates context with pgx.Tx
func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func logQuery(ctx context.Context, q db.Query, args ...any) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}
