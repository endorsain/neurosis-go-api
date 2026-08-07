package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/endorsain/neurosis-go-api/internal/users"
	"github.com/go-chi/chi/v5"
)

type GetUserByIDUseCase interface {
	Execute(ctx context.Context, id string) (users.UserWithProfile, error)
}

type GetUserByIDHandler struct {
	useCase GetUserByIDUseCase
}

func NewGetUserByIDHandler(useCase GetUserByIDUseCase) *GetUserByIDHandler {
	return &GetUserByIDHandler{useCase: useCase}
}

func (h *GetUserByIDHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeJSONError(w, http.StatusBadRequest, "id is required")
		return
	}

	result, err := h.useCase.Execute(r.Context(), id)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

func (h *GetUserByIDHandler) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *GetUserByIDHandler) writeJSONError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
