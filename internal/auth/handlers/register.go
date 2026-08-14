package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

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

func (h *RegisterHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req users.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httptransport.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		httptransport.WriteError(w, http.StatusBadRequest, "username, email and password are required")
		return
	}

	createdUser, err := h.usecase.Execute(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrInvalidInput):
			httptransport.WriteError(w, http.StatusBadRequest, "invalid user input")
		case errors.Is(err, users.ErrUsernameTaken), errors.Is(err, users.ErrEmailTaken):
			httptransport.WriteError(w, http.StatusConflict, "username or email already exists")
		default:
			log.Printf("register user failed: %v", err)
			httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	httptransport.WriteJSON(w, http.StatusCreated, users.UserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
	})
}
