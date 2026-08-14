package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/endorsain/neurosis-go-api/internal/config"
	httptransport "github.com/endorsain/neurosis-go-api/internal/transport/http"
)

type LogoutUseCase interface {
	Execute(ctx context.Context, refreshToken string) error
}

type LogoutHandler struct {
	useCase      LogoutUseCase
	cookieConfig config.RefreshCookieConfig
}

func NewLogoutHandler(useCase LogoutUseCase, cookieConfig config.RefreshCookieConfig) *LogoutHandler {
	return &LogoutHandler{useCase: useCase, cookieConfig: cookieConfig}
}

func (h *LogoutHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieConfig.Name)
	if err == nil {
		if err := h.useCase.Execute(r.Context(), cookie.Value); err != nil {
			log.Printf("logout failed: %v", err)
			httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	} else if err != http.ErrNoCookie {
		log.Printf("read refresh token cookie failed: %v", err)
		httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httptransport.ClearRefreshTokenCookie(w, h.cookieConfig)
	w.WriteHeader(http.StatusNoContent)
}
