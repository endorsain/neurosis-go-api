package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/endorsain/neurosis-go-api/internal/users"
)

type GetCurrentUserUseCase interface {
	Execute(ctx context.Context, id string) (users.UserWithProfile, error)
}

type GetCurrentUserHandler struct {
	useCase GetCurrentUserUseCase
}

func NewGetCurrentUserHandler(useCase GetCurrentUserUseCase) *GetCurrentUserHandler {
	return &GetCurrentUserHandler{useCase: useCase}
}

func (h *GetCurrentUserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		userID = "current-user"
	}

	result, err := h.useCase.Execute(r.Context(), userID)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

func (h *GetCurrentUserHandler) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *GetCurrentUserHandler) writeJSONError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
