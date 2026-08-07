package usecases

import (
	"context"

	"github.com/endorsain/neurosis-go-api/internal/users"
)

type ListUsersUseCase struct {
	userRepository users.UserRepository
}

func NewListUsersUseCase(userRepository users.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{userRepository: userRepository}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context) ([]users.UserSummary, error) {
	return uc.userRepository.List(ctx)
}
