package mailer

import (
	"context"
	"log"
	"os"
	"time"

	logmailer "github.com/mwritter/giftexchanger/services/api/internal/mailer/LogMailer"
	resendmailer "github.com/mwritter/giftexchanger/services/api/internal/mailer/ResendMailer"
)

type Mailer interface {
	SendLoginEmail(ctx context.Context, to, url, code string) error
	SendInvite(ctx context.Context, to, exchangeName, url string) error
}

func New(loginTTL time.Duration) Mailer {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Print("RESEND_API_KEY is not set; emails will be written to the log")
		return logmailer.LogMailer{}
	}
	from := os.Getenv("MAIL_FROM")
	if from == "" {
		log.Fatal("MAIL_FROM is required when RESEND_API_KEY is set")
	}
	return resendmailer.ResendMailer{APIKey: apiKey, From: from, LoginTTL: loginTTL}
}
