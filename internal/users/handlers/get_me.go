package handlers

import (
	"context"
	"log"
	"net/http"

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

// TODO: Esta mal el id deberia llegar de el contexto, gracias a un middleware que valide el accesst_token.
func (h *GetCurrentUserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		userID = "current-user"
	}

	result, err := h.useCase.Execute(r.Context(), userID)
	if err != nil {
		log.Printf("get current user failed: %v", err)
		httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httptransport.WriteJSON(w, http.StatusOK, result)
}
