package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

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
		h.writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		h.writeJSONError(w, http.StatusBadRequest, "username, email and password are required")
		return
	}

	createdUser, err := h.usecase.Execute(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrInvalidInput):
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, users.ErrUsernameTaken), errors.Is(err, users.ErrEmailTaken):
			h.writeJSONError(w, http.StatusConflict, err.Error())
		default:
			h.writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(users.UserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
	})
}

func (h *RegisterHandler) writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
