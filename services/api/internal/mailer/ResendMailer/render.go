package resendmailer

import (
	"bytes"
	"embed"
	"fmt"
	htmltemplate "html/template"
	texttemplate "text/template"
	"time"
)

//go:embed templates/*
var templateFS embed.FS

var (
	htmlTemplates = htmltemplate.Must(htmltemplate.ParseFS(templateFS, "templates/*.html"))
	textTemplates = texttemplate.Must(texttemplate.ParseFS(templateFS, "templates/*.txt"))
)

type renderedEmail struct {
	Subject string
	Text    string
	HTML    string
}

type loginData struct {
	Code   string
	URL    string
	Expiry string
}

type inviteData struct {
	ExchangeName string
	URL          string
}

func renderLogin(url, code string, ttl time.Duration) (renderedEmail, error) {
	data := loginData{Code: code, URL: url, Expiry: humanizeDuration(ttl)}
	return render("login", "Your GiftExchanger login code: "+code, data)
}

func renderInvite(exchangeName, url string) (renderedEmail, error) {
	data := inviteData{ExchangeName: exchangeName, URL: url}
	return render("invite", "You're invited to "+exchangeName, data)
}

func render(name, subject string, data any) (renderedEmail, error) {
	var text, html bytes.Buffer
	if err := textTemplates.ExecuteTemplate(&text, name+".txt", data); err != nil {
		return renderedEmail{}, fmt.Errorf("render %s text: %w", name, err)
	}
	if err := htmlTemplates.ExecuteTemplate(&html, name+".html", data); err != nil {
		return renderedEmail{}, fmt.Errorf("render %s html: %w", name, err)
	}
	return renderedEmail{Subject: subject, Text: text.String(), HTML: html.String()}, nil
}

func humanizeDuration(d time.Duration) string {
	if d <= 0 {
		d = 15 * time.Minute
	}
	if d%time.Hour == 0 {
		return plural(int(d/time.Hour), "hour")
	}
	return plural(int(d.Round(time.Minute)/time.Minute), "minute")
}

func plural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return fmt.Sprintf("%d %ss", n, unit)
}
