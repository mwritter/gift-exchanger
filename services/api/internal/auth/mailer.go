package auth

import (
	"context"
	"log"
)

type Mailer interface {
	SendLoginEmail(ctx context.Context, to, url, code string) error
	SendInvite(ctx context.Context, to, exchangeName, url string) error
}

// LogMailer writes the login URL to the process log. Used in local development
// so we do not need a real inbox.
type LogMailer struct {
	Logger *log.Logger
}

func (m LogMailer) logger() *log.Logger {
	if m.Logger == nil {
		return log.Default()
	}
	return m.Logger
}

func (m LogMailer) SendLoginEmail(_ context.Context, to, url, code string) error {
	m.logger().Printf("login for %s: code %s or link %s", to, code, url)
	return nil
}

func (m LogMailer) SendInvite(_ context.Context, to, exchangeName, url string) error {
	m.logger().Printf("invite for %s to %q: %s", to, exchangeName, url)
	return nil
}
