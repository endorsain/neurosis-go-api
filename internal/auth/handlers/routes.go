package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, registerHandler *RegisterHandler, loginHandler *LoginHandler) {
	r.Post("/auth/register", registerHandler.RegisterUser)
	r.Post("/auth/login", loginHandler.Login)
}
