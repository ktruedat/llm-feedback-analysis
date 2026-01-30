package user

import (
	"context"
	"fmt"
	"time"

	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/user/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/user"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) Create(
	ctx context.Context,
	u *user.User,
	opts ...repository.RepoOption[apprepo.Options],
) error {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	// Convert roles to []string
	roles := u.Roles()
	roleStrings := make([]string, len(roles))
	for i, role := range roles {
		roleStrings[i] = role.String()
	}

	var deletedAt *time.Time
	if u.DeletedAt().IsSome() {
		dt := u.DeletedAt().Unwrap()
		deletedAt = &dt
	}

	_, err := queries.CreateUser(
		ctx, sqlc.CreateUserParams{
			ID:           u.ID(),
			Email:        u.Email().Value(),
			PasswordHash: u.PasswordHash().Value(),
			Roles:        roleStrings,
			Status:       u.Status().String(),
			CreatedAt:    u.CreatedAt(),
			UpdatedAt:    u.UpdatedAt(),
			DeletedAt:    deletedAt,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
