package migrations

import (
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func ApplyMigrations(pool *pgxpool.Pool, logger tracelog.TraceLogger) error {
	logger.Info("applying migrations...")
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect for migrations: %w", err)
	}

	if err := goose.Up(stdlib.OpenDBFromPool(pool), "."); err != nil {
		return fmt.Errorf("failed to run goose up: %w", err)
	}

	logger.Info("applied migrations")
	return nil
}
