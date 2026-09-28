package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/mwritter/giftexchanger/services/api/internal/apitypes"
	"github.com/mwritter/giftexchanger/services/api/internal/auth"
)

// Magic link allows user to login via a link sent to their email
// with a one time authentication token
func (s *Server) requestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req apitypes.MagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err := s.auth.RequestMagicLink(r.Context(), req.Email)
	if errors.Is(err, auth.ErrInvalidEmail) {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if err != nil {
		log.Printf("request magic link: %v", err)
		writeError(w, http.StatusInternalServerError, "could not send magic link")
		return
	}

	writeJSON(w, http.StatusOK, apitypes.StatusResponse{Status: "ok"})
}

func (s *Server) verifyLoginCode(w http.ResponseWriter, r *http.Request) {
	var req apitypes.VerifyLoginCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	sessionRaw, expiresAt, err := s.auth.VerifyLoginCode(r.Context(), req.Email, req.Code)
	if errors.Is(err, auth.ErrInvalidEmail) {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if errors.Is(err, auth.ErrInvalidCode) {
		writeError(w, http.StatusBadRequest, "That code is invalid or has expired.")
		return
	}
	if err != nil {
		log.Printf("verify login code: %v", err)
		writeError(w, http.StatusInternalServerError, "could not log in")
		return
	}

	s.setSessionCookie(w, sessionRaw, expiresAt)
	writeJSON(w, http.StatusOK, apitypes.StatusResponse{Status: "ok"})
}

// Auth callback sets the authorization token in the users browser
// then redirects to the app entry point
func (s *Server) authCallback(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	sessionRaw, expiresAt, err := s.auth.ConsumeMagicLink(r.Context(), token)
	if errors.Is(err, auth.ErrInvalidToken) {
		s.redirectToError(w, r, "invalid_link")
		return
	}
	if err != nil {
		log.Printf("consume magic link: %v", err)
		s.redirectToError(w, r, "unexpected")
		return
	}

	s.setSessionCookie(w, sessionRaw, expiresAt)
	http.Redirect(w, r, s.auth.EntryURL(), http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(s.auth.CookieName())
	raw := ""
	if cookie != nil {
		raw = cookie.Value
	}
	if err := s.auth.Logout(r.Context(), raw); err != nil {
		log.Printf("logout: %v", err)
		writeError(w, http.StatusInternalServerError, "could not log out")
		return
	}
	s.clearSessionCookie(w)
	writeJSON(w, http.StatusOK, apitypes.StatusResponse{Status: "ok"})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, http.StatusOK, apitypes.User{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	})
}

// Safeguard middleware for routes
func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(s.auth.CookieName())
		if err != nil || cookie.Value == "" {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		user, err := s.auth.UserFromSession(r.Context(), cookie.Value)
		if errors.Is(err, auth.ErrUnauthorized) {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		if err != nil {
			log.Printf("load session: %v", err)
			writeError(w, http.StatusInternalServerError, "could not load session")
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
	})
}

// The callback is opened from an email, so failures belong on a page the person
// can read, not in a JSON body.
func (s *Server) redirectToError(w http.ResponseWriter, r *http.Request, reason string) {
	target, err := url.Parse(s.auth.ErrorURL())
	if err != nil {
		log.Printf("parse error url %q: %v", s.auth.ErrorURL(), err)
		writeError(w, http.StatusInternalServerError, "could not complete login")
		return
	}

	query := target.Query()
	query.Set("reason", reason)
	target.RawQuery = query.Encode()

	http.Redirect(w, r, target.String(), http.StatusSeeOther)
}

func (s *Server) setSessionCookie(w http.ResponseWriter, raw string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.auth.CookieName(),
		Value:    raw,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   s.auth.CookieSecure(),
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Server) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.auth.CookieName(),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.auth.CookieSecure(),
		SameSite: http.SameSiteLaxMode,
	})
}
