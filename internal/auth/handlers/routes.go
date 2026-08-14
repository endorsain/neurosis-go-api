package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, registerHandler *RegisterHandler, loginHandler *LoginHandler, refreshHandler *RefreshHandler, logoutHandler *LogoutHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", registerHandler.RegisterUser)
		r.Post("/login", loginHandler.Login)
		r.Post("/refresh", refreshHandler.Refresh)
		r.Post("/logout", logoutHandler.Logout)
	})
}
