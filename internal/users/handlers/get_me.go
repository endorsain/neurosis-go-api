package handlers

import (
	"context"
	"log"
	"net/http"

	authmiddleware "github.com/endorsain/neurosis-go-api/internal/auth/middleware"
	httptransport "github.com/endorsain/neurosis-go-api/internal/transport/http"
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
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok {
		httptransport.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.useCase.Execute(r.Context(), userID)
	if err != nil {
		log.Printf("get current user failed: %v", err)
		httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httptransport.WriteJSON(w, http.StatusOK, result)
}
