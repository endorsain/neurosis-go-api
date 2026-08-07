package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *RegisterHandler) {
	r.Post("/auth", h.RegisterUser)
}
