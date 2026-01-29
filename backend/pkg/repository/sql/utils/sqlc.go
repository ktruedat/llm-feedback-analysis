package utils

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/querier"
)

// QuerierAdapter adapts a querier.PgxQuerier to implement sqlc's DBTX interface.
// This adapter can be used with any sqlc-generated code that requires a DBTX interface.
type QuerierAdapter struct {
	q querier.PgxQuerier
}

// NewQuerierAdapter creates a new adapter that wraps a PgxQuerier to implement DBTX.
func NewQuerierAdapter(q querier.PgxQuerier) *QuerierAdapter {
	return &QuerierAdapter{q: q}
}

// Exec executes a query with the given arguments.
func (a *QuerierAdapter) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return a.q.Exec(ctx, query, args...)
}

// Query executes a query that returns rows.
func (a *QuerierAdapter) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return a.q.Query(ctx, query, args...)
}

// QueryRow executes a query that returns a single row.
func (a *QuerierAdapter) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return a.q.QueryRow(ctx, query, args...)
}

// CopyFrom copies data from a source into a database table.
func (a *QuerierAdapter) CopyFrom(
	ctx context.Context,
	tableName pgx.Identifier,
	columnNames []string,
	rowSrc pgx.CopyFromSource,
) (int64, error) {
	return a.q.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func TimeToPgTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

func OptionalToPgTimestamp(opt optional.Optional[time.Time]) pgtype.Timestamp {
	if opt.IsNone() {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{
		Time:  opt.Unwrap(),
		Valid: true,
	}
}

func PgTimestampToTime(pgTS pgtype.Timestamp) time.Time {
	if !pgTS.Valid {
		return time.Time{}
	}
	return pgTS.Time
}

func PgTimestampToOptional(pgTS pgtype.Timestamp) optional.Optional[time.Time] {
	if !pgTS.Valid {
		return optional.None[time.Time]()
	}
	return optional.Some(pgTS.Time)
}
