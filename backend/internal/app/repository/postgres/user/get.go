package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/user"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) GetByID(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repository.RepoOption[apprepo.Options],
) (*user.User, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	sqlcUser, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return mapSQLCUserToDomain(sqlcUser), nil
}

func (r *repo) GetByEmail(
	ctx context.Context,
	email string,
	opts ...repository.RepoOption[apprepo.Options],
) (*user.User, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	sqlcUser, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return mapSQLCUserToDomain(sqlcUser), nil
}
