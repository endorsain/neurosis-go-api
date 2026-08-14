package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	authDomain "github.com/endorsain/neurosis-go-api/internal/auth"
	"github.com/endorsain/neurosis-go-api/internal/auth/usecases"
	"github.com/endorsain/neurosis-go-api/internal/config"
	httptransport "github.com/endorsain/neurosis-go-api/internal/transport/http"
)

type RefreshTokenUseCase interface {
	Execute(ctx context.Context, refreshToken string) (usecases.RefreshResult, error)
}

type RefreshHandler struct {
	useCase      RefreshTokenUseCase
	cookieConfig config.RefreshCookieConfig
}

type refreshResponse struct {
	AccessToken string `json:"access_token"`
}

func NewRefreshHandler(useCase RefreshTokenUseCase, cookieConfig config.RefreshCookieConfig) *RefreshHandler {
	return &RefreshHandler{useCase: useCase, cookieConfig: cookieConfig}
}

func (h *RefreshHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieConfig.Name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			httptransport.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}

		log.Printf("read refresh token cookie failed: %v", err)
		httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	result, err := h.useCase.Execute(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, authDomain.ErrInvalidRefreshToken) {
			httptransport.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}

		log.Printf("refresh token failed: %v", err)
		httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httptransport.SetRefreshTokenCookie(w, result.RefreshToken, result.RefreshExpiresAt, h.cookieConfig)
	httptransport.WriteJSON(w, http.StatusOK, refreshResponse{AccessToken: result.AccessToken})
}
