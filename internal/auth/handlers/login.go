package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	authDomain "github.com/endorsain/neurosis-go-api/internal/auth"
	"github.com/endorsain/neurosis-go-api/internal/auth/usecases"
	"github.com/endorsain/neurosis-go-api/internal/config"
	httptransport "github.com/endorsain/neurosis-go-api/internal/transport/http"
)

type LoginUseCase interface {
	Execute(ctx context.Context, username, password string) (usecases.LoginResult, error)
}

type LoginHandler struct {
	useCase      LoginUseCase
	cookieConfig config.RefreshCookieConfig
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func NewLoginHandler(useCase LoginUseCase, cookieConfig config.RefreshCookieConfig) *LoginHandler {
	return &LoginHandler{useCase: useCase, cookieConfig: cookieConfig}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httptransport.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Username == "" || req.Password == "" {
		httptransport.WriteError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	result, err := h.useCase.Execute(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, authDomain.ErrInvalidCredentials) {
			httptransport.WriteError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		log.Printf("login failed: %v", err)
		httptransport.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httptransport.SetRefreshTokenCookie(w, result.RefreshToken, result.RefreshExpiresAt, h.cookieConfig)
	httptransport.WriteJSON(w, http.StatusOK, loginResponse{AccessToken: result.AccessToken})
}
