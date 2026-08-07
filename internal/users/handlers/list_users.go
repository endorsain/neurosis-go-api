package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/endorsain/neurosis-go-api/internal/users"
)

type ListUsersUseCase interface {
	Execute(ctx context.Context) ([]users.UserSummary, error)
}

type ListUsersHandler struct {
	useCase ListUsersUseCase
}

func NewListUsersHandler(useCase ListUsersUseCase) *ListUsersHandler {
	return &ListUsersHandler{useCase: useCase}
}

func (h *ListUsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	result, err := h.useCase.Execute(r.Context())
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

func (h *ListUsersHandler) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *ListUsersHandler) writeJSONError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
