package handlers

import (
	"context"
	"log"
	"net/http"

	httptransport "github.com/endorsain/neurosis-go-api/internal/transport/http"
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
		httptransport.WriteError(w, http.StatusBadRequest, "id is required")
		return
	}

	result, err := h.useCase.Execute(r.Context(), id)
	if err != nil {
		log.Printf("get user by id failed (id=%s): %v", id, err)
		httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httptransport.WriteJSON(w, http.StatusOK, result)
}
