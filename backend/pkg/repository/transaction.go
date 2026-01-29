package repository

import (
	"context"
)

type Executor interface {
	IsExecutor()
}

type Transaction interface {
	Executor
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Transactor interface {
	NewTransaction(ctx context.Context, opts ...TransactionOption) (Transaction, error)
}

type TransactionOption func(*TransactionOptions)

type TransactionOptions struct {
	IsolationLevel IsolationLevel
}

type IsolationLevel int

const (
	IsolationLevelReadCommitted IsolationLevel = iota
	IsolationLevelRepeatableRead
	IsolationLevelSerializable
)

func WithIsolationLevel(level IsolationLevel) TransactionOption {
	return func(to *TransactionOptions) {
		to.IsolationLevel = level
	}
}
