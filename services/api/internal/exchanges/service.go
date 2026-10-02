package exchanges

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mwritter/giftexchanger/services/api/internal/apitypes"
	"github.com/mwritter/giftexchanger/services/api/internal/auth"
	"github.com/mwritter/giftexchanger/services/api/internal/mailer"
)

var (
	ErrInvalidName   = errors.New("invalid name")
	ErrInvalidDate   = errors.New("invalid date")
	ErrInvalidBudget = errors.New("invalid budget")
	ErrInvalidEmail  = errors.New("invalid email")
	ErrNotFound      = errors.New("exchange not found")
	ErrForbidden     = errors.New("not the organizer")
)

type Service struct {
	pool    *pgxpool.Pool
	mailer  mailer.Mailer
	baseURL string
}

func NewService(pool *pgxpool.Pool, mailer mailer.Mailer, baseURL string) *Service {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	return &Service{pool: pool, mailer: mailer, baseURL: baseURL}
}

func (s *Service) GetExchanges(ctx context.Context, userID string) ([]apitypes.Exchange, error) {
	rows, err := s.pool.Query(ctx, `
	SELECT
		e.id::text,
		e.organizer_id::text,
		e.name,
		e.exchange_description,
		COALESCE(to_char(e.exchange_date, 'YYYY-MM-DD'), ''),
		e.exchange_budget_cents,
		e.state
	FROM exchanges e
	JOIN exchange_members mine
		ON mine.exchange_id = e.id AND mine.user_id = $1::uuid
	ORDER BY e.created_at DESC
`, userID)

	if err != nil {
		return nil, fmt.Errorf("list exchanges: %w", err)
	}
	defer rows.Close()

	exchanges := []apitypes.Exchange{}
	ids := []string{}
	for rows.Next() {
		var ex apitypes.Exchange
		if err := rows.Scan(
			&ex.ID, &ex.OrganizerID, &ex.Name, &ex.Description,
			&ex.ExchangeDate, &ex.BudgetCents, &ex.State,
		); err != nil {
			return nil, fmt.Errorf("scan exchange: %w", err)
		}
		exchanges = append(exchanges, ex)
		ids = append(ids, ex.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list exchanges: %w", err)
	}
	rows.Close()

	byExchange, err := s.membersByExchange(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range exchanges {
		exchanges[i].Members = withOrganizer(byExchange[exchanges[i].ID], exchanges[i].OrganizerID)
	}
	return exchanges, nil
}

func (s *Service) GetExchangeById(ctx context.Context, userID string, exchangeID string) (apitypes.Exchange, error) {
	var ex apitypes.Exchange
	err := s.pool.QueryRow(ctx, `
	SELECT
		e.id::text,
		e.organizer_id::text,
		e.name,
		e.exchange_description,
		COALESCE(to_char(e.exchange_date, 'YYYY-MM-DD'), ''),
		e.exchange_budget_cents,
		e.state
	FROM exchanges e
	JOIN exchange_members mine
		ON mine.exchange_id = e.id AND mine.user_id = $1::uuid
	WHERE e.id::text = $2
`, userID, exchangeID).Scan(
		&ex.ID, &ex.OrganizerID, &ex.Name, &ex.Description,
		&ex.ExchangeDate, &ex.BudgetCents, &ex.State,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return apitypes.Exchange{}, ErrNotFound
	}
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("get exchange: %w", err)
	}

	byExchange, err := s.membersByExchange(ctx, []string{ex.ID})
	if err != nil {
		return apitypes.Exchange{}, err
	}
	ex.Members = withOrganizer(byExchange[ex.ID], ex.OrganizerID)
	return ex, nil
}

// GetExchangeInvites returns the still-open invites for an exchange. Only the
// organizer may read them, so members get ErrForbidden rather than an empty list.
func (s *Service) GetExchangeInvites(ctx context.Context, userID, exchangeID string) ([]string, error) {
	var organizerID string
	err := s.pool.QueryRow(ctx, `
		SELECT organizer_id::text FROM exchanges WHERE id::text = $1
	`, exchangeID).Scan(&organizerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get exchange organizer: %w", err)
	}
	if organizerID != userID {
		return nil, ErrForbidden
	}

	rows, err := s.pool.Query(ctx, `
		SELECT email
		FROM exchange_invites
		WHERE exchange_id::text = $1
			AND accepted_at IS NULL
			AND declined_at IS NULL
		ORDER BY lower(email)
	`, exchangeID)
	if err != nil {
		return nil, fmt.Errorf("list invites: %w", err)
	}
	defer rows.Close()

	inviteEmails := []string{}
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, fmt.Errorf("scan invite: %w", err)
		}
		inviteEmails = append(inviteEmails, email)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list invites: %w", err)
	}
	return inviteEmails, nil
}

// membersByExchange looks up the members of several exchanges at once so that
// listing exchanges stays a fixed number of queries. IsOrganizer is left to the
// caller, which knows the organizer of each exchange.
func (s *Service) membersByExchange(ctx context.Context, exchangeIDs []string) (map[string][]apitypes.ExchangeMember, error) {
	byExchange := map[string][]apitypes.ExchangeMember{}
	if len(exchangeIDs) == 0 {
		return byExchange, nil
	}

	rows, err := s.pool.Query(ctx, `
		SELECT m.exchange_id::text, u.id::text, u.email, u.display_name
		FROM exchange_members m
		JOIN users u ON u.id = m.user_id
		WHERE m.exchange_id::text = ANY($1::text[])
		ORDER BY m.created_at, lower(u.email)
	`, exchangeIDs)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var exchangeID string
		var member apitypes.ExchangeMember
		if err := rows.Scan(&exchangeID, &member.UserID, &member.Email, &member.DisplayName); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		byExchange[exchangeID] = append(byExchange[exchangeID], member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return byExchange, nil
}

// withOrganizer flags the organizer and guarantees a non-nil slice so the field
// serializes as [] rather than null.
func withOrganizer(members []apitypes.ExchangeMember, organizerID string) []apitypes.ExchangeMember {
	flagged := make([]apitypes.ExchangeMember, 0, len(members))
	for _, member := range members {
		member.IsOrganizer = member.UserID == organizerID
		flagged = append(flagged, member)
	}
	return flagged
}

func (s *Service) CreateExchange(ctx context.Context, organizerID, organizerEmail string, req apitypes.CreateExchangeRequest) (apitypes.Exchange, error) {
	name := strings.TrimSpace(req.Name)
	if len([]rune(name)) < 2 {
		return apitypes.Exchange{}, ErrInvalidName
	}
	description := strings.TrimSpace(req.Description)
	date, err := time.Parse("2006-01-02", strings.TrimSpace(req.ExchangeDate))
	if err != nil {
		return apitypes.Exchange{}, ErrInvalidDate
	}
	if req.BudgetCents != nil && *req.BudgetCents < 0 {
		return apitypes.Exchange{}, ErrInvalidBudget
	}

	organizerEmail, _ = auth.NormalizeEmail(organizerEmail)
	inviteEmails, err := normalizeInviteEmails(req.InviteEmails, organizerEmail)
	if err != nil {
		return apitypes.Exchange{}, err
	}

	token, err := newInviteToken()
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("generate invite token: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx)

	var id string
	err = tx.QueryRow(ctx, `
		INSERT INTO exchanges (
			organizer_id, name, exchange_date, state, invite_token,
			exchange_description, exchange_budget_cents
		)
		VALUES ($1::uuid, $2, $3::date, $4, $5, $6, $7)
		RETURNING id::text
	`, organizerID, name, date.Format("2006-01-02"), apitypes.ExchangeStateDraft, token, description, req.BudgetCents).Scan(&id)
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("insert exchange: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO exchange_members (exchange_id, user_id)
		VALUES ($1::uuid, $2::uuid)
	`, id, organizerID); err != nil {
		return apitypes.Exchange{}, fmt.Errorf("insert organizer membership: %w", err)
	}

	for _, email := range inviteEmails {
		if _, err := tx.Exec(ctx, `
			INSERT INTO exchange_invites (exchange_id, email)
			VALUES ($1::uuid, $2)
		`, id, email); err != nil {
			return apitypes.Exchange{}, fmt.Errorf("insert invite: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return apitypes.Exchange{}, fmt.Errorf("commit: %w", err)
	}

	inviteURL := s.baseURL + "/dashboard"
	for _, email := range inviteEmails {
		if err := s.mailer.SendInvite(ctx, email, name, inviteURL); err != nil {
			log.Printf("send invite to %s: %v", email, err)
		}
	}

	return s.GetExchangeById(ctx, organizerID, id)
}

func (s *Service) UpdateExchange(ctx context.Context, exchangeID, organizerID, organizerEmail string, req apitypes.CreateExchangeRequest) (apitypes.Exchange, error) {
	name := strings.TrimSpace(req.Name)
	if len([]rune(name)) < 2 {
		return apitypes.Exchange{}, ErrInvalidName
	}
	description := strings.TrimSpace(req.Description)
	date, err := time.Parse("2006-01-02", strings.TrimSpace(req.ExchangeDate))
	if err != nil {
		return apitypes.Exchange{}, ErrInvalidDate
	}
	if req.BudgetCents != nil && *req.BudgetCents < 0 {
		return apitypes.Exchange{}, ErrInvalidBudget
	}

	organizerEmail, _ = auth.NormalizeEmail(organizerEmail)
	inviteEmails, err := normalizeInviteEmails(req.InviteEmails, organizerEmail)
	if err != nil {
		return apitypes.Exchange{}, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		UPDATE exchanges
		SET name = $3,
			exchange_description = $4,
			exchange_date = $5::date,
			exchange_budget_cents = $6,
			updated_at = now()
		WHERE id::text = $1 AND organizer_id::text = $2
	`, exchangeID, organizerID, name, description, date.Format("2006-01-02"), req.BudgetCents)
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("update exchange: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apitypes.Exchange{}, ErrNotFound
	}

	rows, err := tx.Query(ctx, `
		SELECT lower(email), accepted_at IS NOT NULL, declined_at IS NOT NULL
		FROM exchange_invites
		WHERE exchange_id::text = $1
	`, exchangeID)
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("list invites: %w", err)
	}

	type inviteState struct {
		accepted bool
		declined bool
	}
	existing := map[string]inviteState{}
	for rows.Next() {
		var email string
		var state inviteState
		if err := rows.Scan(&email, &state.accepted, &state.declined); err != nil {
			rows.Close()
			return apitypes.Exchange{}, fmt.Errorf("scan invite: %w", err)
		}
		existing[email] = state
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return apitypes.Exchange{}, fmt.Errorf("list invites: %w", err)
	}
	rows.Close()

	wanted := make(map[string]struct{}, len(inviteEmails))
	var toSend []string
	for _, email := range inviteEmails {
		wanted[email] = struct{}{}
		state, ok := existing[email]
		if !ok {
			if _, err := tx.Exec(ctx, `
				INSERT INTO exchange_invites (exchange_id, email)
				VALUES ($1::uuid, $2)
			`, exchangeID, email); err != nil {
				return apitypes.Exchange{}, fmt.Errorf("insert invite: %w", err)
			}
			toSend = append(toSend, email)
			continue
		}
		if state.accepted || !state.declined {
			continue
		}
		if _, err := tx.Exec(ctx, `
			UPDATE exchange_invites
			SET declined_at = NULL
			WHERE exchange_id::text = $1 AND lower(email) = $2
		`, exchangeID, email); err != nil {
			return apitypes.Exchange{}, fmt.Errorf("reopen invite: %w", err)
		}
		toSend = append(toSend, email)
	}

	for email, state := range existing {
		if _, ok := wanted[email]; ok || state.accepted || state.declined {
			continue
		}
		if _, err := tx.Exec(ctx, `
			DELETE FROM exchange_invites
			WHERE exchange_id::text = $1
				AND lower(email) = $2
				AND accepted_at IS NULL
				AND declined_at IS NULL
		`, exchangeID, email); err != nil {
			return apitypes.Exchange{}, fmt.Errorf("delete invite: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return apitypes.Exchange{}, fmt.Errorf("commit: %w", err)
	}

	inviteURL := s.baseURL + "/dashboard"
	for _, email := range toSend {
		if err := s.mailer.SendInvite(ctx, email, name, inviteURL); err != nil {
			log.Printf("send invite to %s: %v", email, err)
		}
	}

	return s.GetExchangeById(ctx, organizerID, exchangeID)
}

func (s *Service) DeleteExchange(ctx context.Context, exchangeID, userID string) error {
	tag, err := s.pool.Exec(ctx, `
		DELETE FROM exchanges
		WHERE id::text = $1 AND organizer_id::text = $2
	`, exchangeID, userID)
	if err != nil {
		return fmt.Errorf("delete exchange: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func normalizeInviteEmails(raw []string, organizerEmail string) ([]string, error) {
	seen := make(map[string]struct{}, len(raw))
	emails := make([]string, 0, len(raw))
	for _, rawEmail := range raw {
		email, err := auth.NormalizeEmail(rawEmail)
		if err != nil {
			return nil, ErrInvalidEmail
		}
		if email == organizerEmail {
			continue
		}
		if _, ok := seen[email]; ok {
			continue
		}
		seen[email] = struct{}{}
		emails = append(emails, email)
	}
	return emails, nil
}

func newInviteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
