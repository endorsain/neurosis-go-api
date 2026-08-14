package module

import (
	"database/sql"

	"github.com/endorsain/neurosis-go-api/internal/auth"
	"github.com/endorsain/neurosis-go-api/internal/auth/handlers"
	"github.com/endorsain/neurosis-go-api/internal/auth/usecases"
	"github.com/endorsain/neurosis-go-api/internal/config"
	"github.com/endorsain/neurosis-go-api/internal/users"
	"github.com/go-chi/chi/v5"
)

type Module struct {
	registerHandler *handlers.RegisterHandler
	loginHandler    *handlers.LoginHandler
	refreshHandler  *handlers.RefreshHandler
	logoutHandler   *handlers.LogoutHandler
}

func New(db *sql.DB, userRepository users.UserRepository, authConfig config.AuthConfig) (*Module, error) {
	registerUseCase := usecases.NewRegisterUserUseCase(userRepository)
	authRepository := auth.NewPostgresRepository(db)
	tokenService, err := auth.NewTokenService(authConfig)
	if err != nil {
		return nil, err
	}
	loginUseCase := usecases.NewLoginUseCase(authRepository, tokenService)
	refreshTokenUseCase := usecases.NewRefreshTokenUseCase(authRepository, tokenService)
	logoutUseCase := usecases.NewLogoutUseCase(authRepository)

	return &Module{
		registerHandler: handlers.NewHandler(registerUseCase),
		loginHandler:    handlers.NewLoginHandler(loginUseCase, authConfig.RefreshCookie),
		refreshHandler:  handlers.NewRefreshHandler(refreshTokenUseCase, authConfig.RefreshCookie),
		logoutHandler:   handlers.NewLogoutHandler(logoutUseCase, authConfig.RefreshCookie),
	}, nil
}

func (m *Module) RegisterRoutes(r chi.Router) {
	handlers.RegisterRoutes(r, m.registerHandler, m.loginHandler, m.refreshHandler, m.logoutHandler)
}
