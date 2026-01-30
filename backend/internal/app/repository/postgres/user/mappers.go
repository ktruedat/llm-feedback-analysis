package user

import (
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/user/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/user"
)

// mapSQLCUserToDomain maps a SQLC user model to a domain user entity.
func mapSQLCUserToDomain(sqlcUser sqlc.User) *user.User {
	// Build email value object
	email, _ := user.NewEmail(sqlcUser.Email)

	// Build password hash value object
	passwordHash, _ := user.NewPasswordHash(sqlcUser.PasswordHash)

	// Convert roles from []string to []user.Role
	roles := make([]user.Role, len(sqlcUser.Roles))
	for i, roleStr := range sqlcUser.Roles {
		role, _ := user.NewRole(roleStr)
		roles[i] = role
	}

	// Build status value object
	status, _ := user.NewUserStatus(sqlcUser.Status)

	// Build domain entity using builder
	builder := user.NewBuilder().
		WithID(sqlcUser.ID).
		WithEmail(email).
		WithPasswordHash(passwordHash).
		WithRoles(roles).
		WithStatus(status).
		WithCreatedAt(sqlcUser.CreatedAt).
		WithUpdatedAt(sqlcUser.UpdatedAt)

	// Handle deleted_at (nullable timestamp)
	if sqlcUser.DeletedAt != nil {
		builder.WithDeletedAt(*sqlcUser.DeletedAt)
	}

	return builder.BuildUnchecked()
}
