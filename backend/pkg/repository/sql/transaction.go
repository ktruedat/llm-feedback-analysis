package sql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/querier"
)

type transactor struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) repository.Transactor {
	return &transactor{pool: pool}
}

func (t *transactor) NewTransaction(ctx context.Context, opts ...repository.TransactionOption) (
	repository.Transaction,
	error,
) {
	pgxTx, err := t.pool.BeginTx(ctx, toPgxTxOptions(opts))
	if err != nil {
		return nil, fmt.Errorf("transaction could not be started: %w", err)
	}

	return &transaction{pgxTx: pgxTx}, nil
}

func toPgxTxOptions(opts []repository.TransactionOption) pgx.TxOptions {
	var txOptions repository.TransactionOptions
	for _, opt := range opts {
		opt(&txOptions)
	}

	var pgxIsolationLevel pgx.TxIsoLevel
	switch txOptions.IsolationLevel {
	case repository.IsolationLevelSerializable:
		pgxIsolationLevel = pgx.Serializable
	case repository.IsolationLevelReadCommitted:
		pgxIsolationLevel = pgx.ReadCommitted
	case repository.IsolationLevelRepeatableRead:
		pgxIsolationLevel = pgx.RepeatableRead
	default:
		pgxIsolationLevel = pgx.ReadCommitted
	}

	return pgx.TxOptions{
		IsoLevel: pgxIsolationLevel,
	}
}

var (
	_ repository.Transaction = (*transaction)(nil)
	_ querier.PgxQuerier     = (*transaction)(nil)
)

type transaction struct {
	pgxTx pgx.Tx
}

func (t *transaction) Commit(ctx context.Context) error {
	if err := t.pgxTx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction could not be committed: %w", err)
	}

	return nil
}

func (t *transaction) Rollback(ctx context.Context) error {
	if err := t.pgxTx.Rollback(ctx); err != nil {
		return fmt.Errorf("transaction could not be rolled back: %w", err)
	}

	return nil
}

func (*transaction) IsExecutor() {}

func (t *transaction) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return t.pgxTx.Exec(ctx, query, args...)
}

func (t *transaction) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return t.pgxTx.Query(ctx, sql, args...)
}

func (t *transaction) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return t.pgxTx.QueryRow(ctx, sql, args...)
}

func (t *transaction) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	return t.pgxTx.SendBatch(ctx, batch)
}

func (t *transaction) CopyFrom(
	ctx context.Context,
	tableName pgx.Identifier,
	columnNames []string,
	rowSrc pgx.CopyFromSource,
) (int64, error) {
	return t.pgxTx.CopyFrom(ctx, tableName, columnNames, rowSrc)
}
