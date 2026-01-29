package querier

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
)

type PgxQuerier interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (
		int64,
		error,
	)
}

type Scanner interface {
	Scan(dest ...any) error
}

var (
	_ PgxQuerier          = (*PgxPool)(nil)
	_ repository.Executor = (*PgxPool)(nil)
)

type PgxPool struct {
	pool *pgxpool.Pool
}

func NewPgxPool(pool *pgxpool.Pool) *PgxPool {
	return &PgxPool{pool: pool}
}

func (p *PgxPool) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, query, args...)
}

func (p *PgxPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *PgxPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *PgxPool) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	return p.pool.SendBatch(ctx, batch)
}

func (p *PgxPool) CopyFrom(
	ctx context.Context,
	tableName pgx.Identifier,
	columnNames []string,
	rowSrc pgx.CopyFromSource,
) (int64, error) {
	return p.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (*PgxPool) IsExecutor() {}
