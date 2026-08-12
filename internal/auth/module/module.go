package module

import (
	"github.com/endorsain/neurosis-go-api/internal/auth/handlers"
	"github.com/endorsain/neurosis-go-api/internal/auth/usecases"
	"github.com/endorsain/neurosis-go-api/internal/users"
	"github.com/go-chi/chi/v5"
)

type Module struct {
	registerHandler *handlers.RegisterHandler
}

func New(userRepository users.UserRepository) *Module {
	registerUseCase := usecases.NewRegisterUserUseCase(userRepository)

	return &Module{
		registerHandler: handlers.NewHandler(registerUseCase),
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	handlers.RegisterRoutes(r, m.registerHandler)
}
