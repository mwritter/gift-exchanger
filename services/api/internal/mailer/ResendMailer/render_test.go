package resendmailer

import (
	"strings"
	"testing"
	"time"
)

func TestRenderLogin(t *testing.T) {
	url := "https://example.com/api/auth/callback?token=abc&x=1"
	email, err := renderLogin(url, "123456", 15*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if email.Subject != "Your GiftExchanger login code: 123456" {
		t.Errorf("subject = %q", email.Subject)
	}
	for name, body := range map[string]string{"text": email.Text, "html": email.HTML} {
		if !strings.Contains(body, "123456") {
			t.Errorf("%s body missing code", name)
		}
		if !strings.Contains(body, "15 minutes") {
			t.Errorf("%s body missing expiry", name)
		}
	}
	if !strings.Contains(email.Text, url) {
		t.Error("text body missing url")
	}
	if !strings.Contains(email.HTML, "https://example.com/api/auth/callback?token=abc&amp;x=1") {
		t.Error("html body missing url")
	}
}

func TestRenderInviteEscapesExchangeName(t *testing.T) {
	email, err := renderInvite(`<script>alert("x")</script>`, "https://example.com/dashboard")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(email.HTML, "<script>") {
		t.Error("html body contains unescaped exchange name")
	}
	if !strings.Contains(email.HTML, "&lt;script&gt;") {
		t.Error("html body missing escaped exchange name")
	}
	if !strings.Contains(email.Text, "https://example.com/dashboard") || !strings.Contains(email.HTML, "https://example.com/dashboard") {
		t.Error("body missing url")
	}
}

func TestHumanizeDuration(t *testing.T) {
	cases := map[time.Duration]string{
		0:                "15 minutes",
		time.Minute:      "1 minute",
		15 * time.Minute: "15 minutes",
		time.Hour:        "1 hour",
		2 * time.Hour:    "2 hours",
		90 * time.Minute: "90 minutes",
	}
	for d, want := range cases {
		if got := humanizeDuration(d); got != want {
			t.Errorf("humanizeDuration(%v) = %q, want %q", d, got, want)
		}
	}
}
