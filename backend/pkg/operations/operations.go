package operations

import (
	"context"
	"fmt"

	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
)

type TxExecFunc func(ctx context.Context, tx repository.Transaction) error

func RunGenericTransaction(
	ctx context.Context,
	transactor repository.Transactor,
	exec TxExecFunc,
) (err error) {
	tx, err := transactor.NewTransaction(ctx, repository.WithIsolationLevel(repository.IsolationLevelReadCommitted))
	if err != nil {
		return fmt.Errorf("could not init transaction: %w", err)
	}
	defer deferRollbackOnError(ctx, tx, &err)()

	if err = exec(ctx, tx); err != nil {
		return fmt.Errorf("failed to execute transaction logic: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func deferRollbackOnError(ctx context.Context, tx repository.Transaction, errPtr *error) func() {
	return func() {
		if *errPtr != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				*errPtr = fmt.Errorf("failed to rollback transaction: %w: %v", rbErr, *errPtr)
			}
		}
	}
}
