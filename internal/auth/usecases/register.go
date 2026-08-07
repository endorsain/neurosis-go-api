package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/endorsain/neurosis-go-api/internal/users"
)

type RegisterUserUseCase struct {
	userRepository users.UserRepository
}

func NewRegisterUserUseCase(userRepository users.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepository: userRepository}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, username, email, password string) (users.User, error) {

	_, err := uc.userRepository.FindByUsername(ctx, username)
	if err == nil {
		return users.User{}, users.ErrUsernameTaken
	}
	if !errors.Is(err, users.ErrNotFound) {
		return users.User{}, fmt.Errorf("verify username: %w", err)
	}

	_, err = uc.userRepository.FindByEmail(ctx, email)
	if err == nil {
		return users.User{}, users.ErrEmailTaken
	}
	if !errors.Is(err, users.ErrNotFound) {
		return users.User{}, fmt.Errorf("verify email: %w", err)
	}

	user := users.User{
		Username:     username,
		Email:        email,
		PasswordHash: password,
	}
	profile := users.UserProfile{}

	createdUser, err := uc.userRepository.Create(ctx, user, profile)
	if err != nil {
		return users.User{}, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}
