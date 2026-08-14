package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, getCurrentUserHandler *GetCurrentUserHandler, getUserByIDHandler *GetUserByIDHandler, listUsersHandler *ListUsersHandler, requireAuthentication func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.With(requireAuthentication).Get("/me", getCurrentUserHandler.GetCurrentUser)
		r.Get("/", listUsersHandler.ListUsers)
		r.Get("/{id}", getUserByIDHandler.GetUserByID)
	})
}
