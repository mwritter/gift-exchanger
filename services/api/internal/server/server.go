package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mwritter/giftexchanger/services/api/internal/auth"
)

type Server struct {
	pool *pgxpool.Pool
	auth *auth.Service
}

func New(pool *pgxpool.Pool, authService *auth.Service) http.Handler {
	s := &Server{pool: pool, auth: authService}
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", s.health)
		r.Get("/ready", s.ready)

		r.Post("/auth/magic-link", s.requestMagicLink)
		r.Post("/auth/verify-code", s.verifyLoginCode)
		r.Get("/auth/callback", s.authCallback)
		r.Post("/auth/logout", s.logout)

		r.Group(func(r chi.Router) {
			r.Use(s.requireAuth)
			r.Get("/me", s.me)
		})
	})

	return r
}
