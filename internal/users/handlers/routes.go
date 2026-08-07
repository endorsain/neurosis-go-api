package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, getCurrentUserHandler *GetCurrentUserHandler, getUserByIDHandler *GetUserByIDHandler, listUsersHandler *ListUsersHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/me", getCurrentUserHandler.GetCurrentUser)
		r.Get("/", listUsersHandler.ListUsers)
		r.Get("/{id}", getUserByIDHandler.GetUserByID)
	})
}
