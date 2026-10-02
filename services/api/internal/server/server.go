package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mwritter/giftexchanger/services/api/internal/auth"
	"github.com/mwritter/giftexchanger/services/api/internal/exchanges"
)

type Server struct {
	pool      *pgxpool.Pool
	auth      *auth.Service
	exchanges *exchanges.Service
}

func New(pool *pgxpool.Pool, authService *auth.Service, exchangesService *exchanges.Service) http.Handler {
	s := &Server{pool: pool, auth: authService, exchanges: exchangesService}
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

			r.Get("/exchanges", s.getExchanges)
			r.Get("/exchanges/{exchangeID}", s.getExchange)
			r.Get("/exchanges/{exchangeID}/invites", s.getExchangeInvites)
			r.Post("/exchanges/{exchangeID}/invites/accept", s.acceptExchangeInvite)
			r.Post("/exchanges/{exchangeID}/invites/decline", s.declineExchangeInvite)
			r.Post("/exchanges", s.createExchange)
			r.Put("/exchanges/{exchangeID}", s.updateExchange)
			r.Delete("/exchanges/{exchangeID}", s.deleteExchange)
		})
	})

	return r
}
