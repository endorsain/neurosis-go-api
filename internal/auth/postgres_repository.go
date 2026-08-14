package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/endorsain/neurosis-go-api/internal/users"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindLoginUser(ctx context.Context, username string) (users.User, error) {
	const query = `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE username = $1
	`

	var user users.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users.User{}, users.ErrNotFound
		}
		return users.User{}, fmt.Errorf("find login user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) FindUserRoles(ctx context.Context, userID string) ([]string, error) {
	const query = `
		SELECT r.name
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find user roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("scan user role: %w", err)
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user roles: %w", err)
	}

	return roles, nil
}

func (r *PostgresRepository) SaveRefreshToken(ctx context.Context, token RefreshToken) error {
	const query = `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`

	if _, err := r.db.ExecContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt); err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}

	return nil
}
