package httptransport

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	httpmiddleware "github.com/endorsain/neurosis-go-api/internal/transport/http/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	UseCommonMiddleware(r)
	return r
}

func UseCommonMiddleware(r *chi.Mux) {
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
}

// Handle registers an error-returning handler, routing any error it returns
// through the centralized HTTP error middleware (see transport/http/middleware).
func Handle(fn httpmiddleware.HandlerFunc) http.HandlerFunc {
	return httpmiddleware.Errors(fn)
}

// SetErrorLogger configures the logger used to record request failures
// handled by Handle. It should be called once, from the composition root,
// before the router starts serving requests.
func SetErrorLogger(logger *slog.Logger) {
	httpmiddleware.SetLogger(logger)
}
