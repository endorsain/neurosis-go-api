package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (User, error) {
	const query = `
		SELECT id, username, email, password_hash
		FROM users
		WHERE username = $1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("find user by username: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	const query = `
		SELECT id, username, email, password_hash
		FROM users
		WHERE email = $1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user User, profile UserProfile) (User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	const insertUserQuery = `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, insertUserQuery, user.Username, user.Email, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return User{}, err
	}

	profile.UserID = user.ID
	const insertProfileQuery = `
		INSERT INTO user_profiles (user_id)
		VALUES ($1)
	`

	_, err = tx.ExecContext(ctx, insertProfileQuery, profile.UserID)
	if err != nil {
		return User{}, err
	}

	if err = tx.Commit(); err != nil {
		return User{}, err
	}

	return user, nil
}

// TODO: MAURO RECORDA MEJORAR LOS LOGS PARA ERRORES.
func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (UserWithProfile, error) {
	const query = `
		SELECT 
		    u.id, u.username, u.email, u.password_hash, u.created_at,
			p.id, p.user_id, COALESCE(p.display_name, ''), COALESCE(p.biography, ''), p.updated_at
		FROM users u
		LEFT JOIN user_profiles p ON p.user_id = u.id
		WHERE u.id = $1
	`

	var user UserWithProfile
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Profile.ID,
		&user.Profile.UserID,
		&user.Profile.DisplayName,
		&user.Profile.Biography,
		&user.Profile.UpdatedAt,
	)
	if err != nil {
		log.Printf("get user by id: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return UserWithProfile{}, ErrNotFound
		}
		return UserWithProfile{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) List(ctx context.Context) ([]UserSummary, error) {
	const query = `
		SELECT id, username
		FROM users
		ORDER BY username ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []UserSummary
	for rows.Next() {
		var user UserSummary
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("scan user summary: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}
