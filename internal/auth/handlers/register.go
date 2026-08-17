package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	apperrors "github.com/endorsain/neurosis-go-api/internal/errors"
	httptransport "github.com/endorsain/neurosis-go-api/internal/transport/http"
	"github.com/endorsain/neurosis-go-api/internal/users"
)

type RegisterUserUseCase interface {
	Execute(ctx context.Context, username, email, password string) (users.User, error)
}

type RegisterHandler struct {
	usecase RegisterUserUseCase
}

func NewHandler(registerUserUseCase RegisterUserUseCase) *RegisterHandler {
	return &RegisterHandler{usecase: registerUserUseCase}
}

func (h *RegisterHandler) RegisterUser(w http.ResponseWriter, r *http.Request) error {
	var req users.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid JSON body", http.StatusBadRequest)
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "username, email and password are required", http.StatusBadRequest)
	}

	createdUser, err := h.usecase.Execute(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	httptransport.WriteJSON(w, http.StatusCreated, users.UserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
	})
	return nil
}
