package resendmailer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestResendMailerSendLoginEmail(t *testing.T) {
	var got resendEmailRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/emails" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if h := r.Header.Get("Authorization"); h != "Bearer re_test" {
			t.Errorf("authorization = %q", h)
		}
		if h := r.Header.Get("Content-Type"); h != "application/json" {
			t.Errorf("content-type = %q", h)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.Write([]byte(`{"id":"email_123"}`))
	}))
	defer srv.Close()

	m := ResendMailer{APIKey: "re_test", From: "GiftExchanger <hello@example.com>", BaseURL: srv.URL}
	if err := m.SendLoginEmail(context.Background(), "a@example.com", "https://example.com/cb", "654321"); err != nil {
		t.Fatal(err)
	}

	if got.From != "GiftExchanger <hello@example.com>" {
		t.Errorf("from = %q", got.From)
	}
	if len(got.To) != 1 || got.To[0] != "a@example.com" {
		t.Errorf("to = %v", got.To)
	}
	if !strings.Contains(got.Subject, "654321") {
		t.Errorf("subject = %q", got.Subject)
	}
	if !strings.Contains(got.HTML, "654321") || !strings.Contains(got.Text, "https://example.com/cb") {
		t.Error("body missing code or url")
	}
}

func TestResendMailerAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"statusCode":403,"name":"validation_error","message":"The domain is not verified."}`))
	}))
	defer srv.Close()

	m := ResendMailer{APIKey: "re_secret", From: "x@example.com", BaseURL: srv.URL}
	err := m.SendInvite(context.Background(), "a@example.com", "Family", "https://example.com/dashboard")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "403") || !strings.Contains(err.Error(), "The domain is not verified.") {
		t.Errorf("error = %v", err)
	}
	if strings.Contains(err.Error(), "re_secret") {
		t.Error("error leaks api key")
	}
}

func TestResendMailerNonJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer srv.Close()

	m := ResendMailer{APIKey: "re_test", From: "x@example.com", BaseURL: srv.URL}
	err := m.SendInvite(context.Background(), "a@example.com", "Family", "https://example.com/dashboard")
	if err == nil || !strings.Contains(err.Error(), "502") {
		t.Errorf("error = %v", err)
	}
}

func TestResendMailerCancelledContext(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	defer srv.Close()
	defer close(release)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	m := ResendMailer{APIKey: "re_test", From: "x@example.com", BaseURL: srv.URL}
	err := m.SendLoginEmail(ctx, "a@example.com", "https://example.com/cb", "111111")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("error = %v, want context deadline exceeded", err)
	}
}
