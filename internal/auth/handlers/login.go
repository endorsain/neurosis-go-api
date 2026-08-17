package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/endorsain/neurosis-go-api/internal/auth/usecases"
	"github.com/endorsain/neurosis-go-api/internal/config"
	apperrors "github.com/endorsain/neurosis-go-api/internal/errors"
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

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid JSON body", http.StatusBadRequest)
	}

	if req.Username == "" || req.Password == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "username and password are required", http.StatusBadRequest)
	}

	result, err := h.useCase.Execute(r.Context(), req.Username, req.Password)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	httptransport.SetRefreshTokenCookie(w, result.RefreshToken, result.RefreshExpiresAt, h.cookieConfig)
	httptransport.WriteJSON(w, http.StatusOK, loginResponse{AccessToken: result.AccessToken})
	return nil
}
