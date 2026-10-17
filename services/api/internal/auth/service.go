package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mwritter/giftexchanger/services/api/internal/db"
	"github.com/mwritter/giftexchanger/services/api/internal/mailer"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidToken = errors.New("invalid or expired magic link")
	ErrInvalidCode  = errors.New("invalid or expired code")
	ErrUnauthorized = errors.New("unauthorized")
)

const maxCodeAttempts = 5

type Config struct {
	BaseURL      string
	EntryURL     string
	ErrorURL     string
	MagicLinkTTL time.Duration
	SessionTTL   time.Duration
	CookieSecure bool
	CookieName   string
}

func (c Config) withDefaults() Config {
	if c.BaseURL == "" {
		c.BaseURL = "http://localhost:3000"
	}
	if c.EntryURL == "" {
		c.EntryURL = strings.TrimRight(c.BaseURL, "/") + "/dashboard"
	}
	if c.ErrorURL == "" {
		c.ErrorURL = strings.TrimRight(c.BaseURL, "/") + "/auth/error"
	}
	if c.MagicLinkTTL == 0 {
		c.MagicLinkTTL = 15 * time.Minute
	}
	if c.SessionTTL == 0 {
		c.SessionTTL = 30 * 24 * time.Hour
	}
	if c.CookieName == "" {
		c.CookieName = "giftexchanger_session"
	}
	return c
}

type User struct {
	ID          string
	Email       string
	DisplayName string
}

type Service struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	mailer  mailer.Mailer
	cfg     Config
}

func NewService(pool *pgxpool.Pool, mailer mailer.Mailer, cfg Config) *Service {
	return &Service{pool: pool, queries: db.New(pool), mailer: mailer, cfg: cfg.withDefaults()}
}

func (s *Service) CookieName() string        { return s.cfg.CookieName }
func (s *Service) BaseURL() string           { return strings.TrimRight(s.cfg.BaseURL, "/") }
func (s *Service) EntryURL() string          { return s.cfg.EntryURL }
func (s *Service) ErrorURL() string          { return s.cfg.ErrorURL }
func (s *Service) CookieSecure() bool        { return s.cfg.CookieSecure }
func (s *Service) SessionTTL() time.Duration { return s.cfg.SessionTTL }

func NormalizeEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if email == "" {
		return "", ErrInvalidEmail
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", ErrInvalidEmail
	}
	return email, nil
}

func (s *Service) RequestMagicLink(ctx context.Context, rawEmail string) error {
	email, err := NormalizeEmail(rawEmail)
	if err != nil {
		return err
	}

	raw, err := newRandomToken()
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}
	hash := HashToken(raw)
	code, err := newLoginCode()
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}
	expires := time.Now().Add(s.cfg.MagicLinkTTL)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := s.queries.WithTx(tx)

	if err := q.InvalidateMagicLinkTokens(ctx, email); err != nil {
		return fmt.Errorf("invalidate prior tokens: %w", err)
	}

	codeHash := HashLoginCode(email, code)
	if err := q.InsertMagicLinkToken(ctx, db.InsertMagicLinkTokenParams{
		Email:     email,
		TokenHash: hash,
		CodeHash:  &codeHash,
		ExpiresAt: expires,
	}); err != nil {
		return fmt.Errorf("insert token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	link := fmt.Sprintf("%s/api/auth/callback?token=%s", s.BaseURL(), raw)
	if err := s.mailer.SendLoginEmail(ctx, email, link, code); err != nil {
		return fmt.Errorf("send login email: %w", err)
	}
	return nil
}

func (s *Service) ConsumeMagicLink(ctx context.Context, rawToken string) (sessionRaw string, expiresAt time.Time, err error) {
	if strings.TrimSpace(rawToken) == "" {
		return "", time.Time{}, ErrInvalidToken
	}

	hash := HashToken(rawToken)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := s.queries.WithTx(tx)

	email, err := q.ConsumeMagicLinkToken(ctx, hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", time.Time{}, ErrInvalidToken
	}
	if err != nil {
		return "", time.Time{}, fmt.Errorf("consume token: %w", err)
	}

	sessionRaw, expiresAt, err = s.startSession(ctx, q, email)
	if err != nil {
		return "", time.Time{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", time.Time{}, fmt.Errorf("commit: %w", err)
	}
	return sessionRaw, expiresAt, nil
}

// VerifyLoginCode is the same-tab alternative to the emailed link. Only the
// newest unconsumed token for the email is checked, and it is consumed after
// maxCodeAttempts wrong guesses so a 6-digit code cannot be brute-forced.
func (s *Service) VerifyLoginCode(ctx context.Context, rawEmail, rawCode string) (sessionRaw string, expiresAt time.Time, err error) {
	email, err := NormalizeEmail(rawEmail)
	if err != nil {
		return "", time.Time{}, err
	}
	code := strings.TrimSpace(rawCode)
	if code == "" {
		return "", time.Time{}, ErrInvalidCode
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := s.queries.WithTx(tx)

	token, err := q.GetActiveLoginCodeForUpdate(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", time.Time{}, ErrInvalidCode
	}
	if err != nil {
		return "", time.Time{}, fmt.Errorf("load code: %w", err)
	}
	var codeHash string
	if token.CodeHash != nil {
		codeHash = *token.CodeHash
	}

	if subtle.ConstantTimeCompare([]byte(codeHash), []byte(HashLoginCode(email, code))) != 1 {
		if err := q.RecordLoginCodeAttempt(ctx, db.RecordLoginCodeAttemptParams{
			ID:          token.ID,
			MaxAttempts: maxCodeAttempts,
		}); err != nil {
			return "", time.Time{}, fmt.Errorf("record attempt: %w", err)
		}
		if err := tx.Commit(ctx); err != nil {
			return "", time.Time{}, fmt.Errorf("commit: %w", err)
		}
		return "", time.Time{}, ErrInvalidCode
	}

	if err := q.ConsumeMagicLinkTokenByID(ctx, token.ID); err != nil {
		return "", time.Time{}, fmt.Errorf("consume code: %w", err)
	}
	sessionRaw, expiresAt, err = s.startSession(ctx, q, email)
	if err != nil {
		return "", time.Time{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", time.Time{}, fmt.Errorf("commit: %w", err)
	}
	return sessionRaw, expiresAt, nil
}

// startSession creates the user on first login and inserts a session, inside
// the caller's transaction.
func (s *Service) startSession(ctx context.Context, q *db.Queries, email string) (sessionRaw string, expiresAt time.Time, err error) {
	userID, err := q.UpsertUserByEmail(ctx, email)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("upsert user: %w", err)
	}

	sessionRaw, err = newRandomToken()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("generate session: %w", err)
	}
	expiresAt = time.Now().Add(s.cfg.SessionTTL)
	if err := q.InsertSession(ctx, db.InsertSessionParams{
		UserID:    userID,
		TokenHash: HashToken(sessionRaw),
		ExpiresAt: expiresAt,
	}); err != nil {
		return "", time.Time{}, fmt.Errorf("insert session: %w", err)
	}
	return sessionRaw, expiresAt, nil
}

func (s *Service) Logout(ctx context.Context, sessionRaw string) error {
	if sessionRaw == "" {
		return nil
	}
	if err := s.queries.DeleteSession(ctx, HashToken(sessionRaw)); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (s *Service) UserFromSession(ctx context.Context, sessionRaw string) (*User, error) {
	if sessionRaw == "" {
		return nil, ErrUnauthorized
	}

	row, err := s.queries.GetUserBySession(ctx, HashToken(sessionRaw))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, fmt.Errorf("lookup session: %w", err)
	}
	return &User{ID: row.ID, Email: row.Email, DisplayName: row.DisplayName}, nil
}
