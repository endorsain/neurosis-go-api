package auth

import (
	"context"
	"time"

	"github.com/endorsain/neurosis-go-api/internal/users"
)

type Repository interface {
	FindLoginUser(ctx context.Context, username string) (users.User, error)
	FindUserRoles(ctx context.Context, userID string) ([]string, error)
	SaveRefreshToken(ctx context.Context, token RefreshToken) error
}

type RefreshToken struct {
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}
