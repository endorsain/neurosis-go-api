package handlers

import (
	"context"
	"log"
	"net/http"

	httptransport "github.com/endorsain/neurosis-go-api/internal/transport/http"
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
		log.Printf("list users failed: %v", err)
		httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httptransport.WriteJSON(w, http.StatusOK, result)
}
