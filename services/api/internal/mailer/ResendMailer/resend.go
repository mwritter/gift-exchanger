package resendmailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultResendBaseURL = "https://api.resend.com"

// ResendMailer sends email through the Resend HTTP API.
type ResendMailer struct {
	APIKey string
	From   string
	// LoginTTL is shown in the login email; it should match the magic link TTL.
	LoginTTL   time.Duration
	BaseURL    string
	HTTPClient *http.Client
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
}

type resendErrorResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (m ResendMailer) SendLoginEmail(ctx context.Context, to, url, code string) error {
	email, err := renderLogin(url, code, m.LoginTTL)
	if err != nil {
		return err
	}
	return m.send(ctx, to, email)
}

func (m ResendMailer) SendInvite(ctx context.Context, to, exchangeName, url string) error {
	email, err := renderInvite(exchangeName, url)
	if err != nil {
		return err
	}
	return m.send(ctx, to, email)
}

func (m ResendMailer) send(ctx context.Context, to string, email renderedEmail) error {
	body, err := json.Marshal(resendEmailRequest{
		From:    m.From,
		To:      []string{to},
		Subject: email.Subject,
		HTML:    email.HTML,
		Text:    email.Text,
	})
	if err != nil {
		return fmt.Errorf("encode resend request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.baseURL()+"/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client().Do(req)
	if err != nil {
		return fmt.Errorf("resend request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		io.Copy(io.Discard, resp.Body)
		return nil
	}

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	var apiErr resendErrorResponse
	if json.Unmarshal(raw, &apiErr) == nil && apiErr.Message != "" {
		return fmt.Errorf("resend: %d %s: %s", resp.StatusCode, apiErr.Name, apiErr.Message)
	}
	return fmt.Errorf("resend: unexpected status %d", resp.StatusCode)
}

func (m ResendMailer) baseURL() string {
	if m.BaseURL == "" {
		return defaultResendBaseURL
	}
	return strings.TrimRight(m.BaseURL, "/")
}

func (m ResendMailer) client() *http.Client {
	if m.HTTPClient == nil {
		return &http.Client{Timeout: 10 * time.Second}
	}
	return m.HTTPClient
}
