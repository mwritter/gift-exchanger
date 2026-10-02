package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mwritter/giftexchanger/services/api/internal/auth"
	"github.com/mwritter/giftexchanger/services/api/internal/exchanges"
	"github.com/mwritter/giftexchanger/services/api/internal/mailer"
	"github.com/mwritter/giftexchanger/services/api/internal/server"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to create database pool: %v", err)
	}
	defer pool.Close()

	mailerService := mailer.LogMailer{}
	baseURL := envOr("APP_BASE_URL", "http://localhost:3000")

	authService := auth.NewService(pool, mailerService, auth.Config{
		BaseURL:      baseURL,
		EntryURL:     envOr("APP_ENTRY_URL", baseURL+"/dashboard"),
		ErrorURL:     envOr("APP_ERROR_URL", baseURL+"/auth/error"),
		MagicLinkTTL: durationOr("MAGIC_LINK_TTL", 15*time.Minute),
		SessionTTL:   durationOr("SESSION_TTL", 30*24*time.Hour),
		CookieSecure: boolOr("COOKIE_SECURE", false),
	})

	exchangesService := exchanges.NewService(pool, mailerService, baseURL)

	handler := server.New(pool, authService, exchangesService)

	port := envOr("PORT", "8080")
	if port[0] != ':' {
		port = ":" + port
	}
	log.Printf("Server starting on port %s...", port)

	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
