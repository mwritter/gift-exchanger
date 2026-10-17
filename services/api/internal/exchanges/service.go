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
	"github.com/mwritter/giftexchanger/services/api/internal/db"
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
	queries *db.Queries
	mailer  mailer.Mailer
	baseURL string
}

func NewService(pool *pgxpool.Pool, mailer mailer.Mailer, baseURL string) *Service {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	return &Service{pool: pool, queries: db.New(pool), mailer: mailer, baseURL: baseURL}
}

func (s *Service) GetExchanges(ctx context.Context, userID string) ([]apitypes.Exchange, error) {
	rows, err := s.queries.ListExchangesForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list exchanges: %w", err)
	}

	exchanges := make([]apitypes.Exchange, 0, len(rows))
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		exchanges = append(exchanges, exchangeFromRow(db.GetExchangeForUserRow(row)))
		ids = append(ids, row.ID)
	}

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
	row, err := s.queries.GetExchangeForUser(ctx, db.GetExchangeForUserParams{
		UserID:     userID,
		ExchangeID: exchangeID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return apitypes.Exchange{}, ErrNotFound
	}
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("get exchange: %w", err)
	}
	ex := exchangeFromRow(row)

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
	organizerID, err := s.queries.GetExchangeOrganizerID(ctx, exchangeID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get exchange organizer: %w", err)
	}
	if organizerID != userID {
		return nil, ErrForbidden
	}

	inviteEmails, err := s.queries.ListOpenInviteEmails(ctx, exchangeID)
	if err != nil {
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

	rows, err := s.queries.ListMembersForExchanges(ctx, exchangeIDs)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	for _, row := range rows {
		byExchange[row.ExchangeID] = append(byExchange[row.ExchangeID], apitypes.ExchangeMember{
			UserID:      row.UserID,
			Email:       row.Email,
			DisplayName: row.DisplayName,
		})
	}
	return byExchange, nil
}

// exchangeFromRow leaves Members unset; callers fill it from membersByExchange.
func exchangeFromRow(row db.GetExchangeForUserRow) apitypes.Exchange {
	return apitypes.Exchange{
		ID:           row.ID,
		OrganizerID:  row.OrganizerID,
		Name:         row.Name,
		Description:  row.ExchangeDescription,
		ExchangeDate: row.ExchangeDate,
		BudgetCents:  intPtr(row.ExchangeBudgetCents),
		State:        apitypes.ExchangeState(row.State),
	}
}

func intPtr(v *int32) *int {
	if v == nil {
		return nil
	}
	n := int(*v)
	return &n
}

func int32Ptr(v *int) *int32 {
	if v == nil {
		return nil
	}
	n := int32(*v)
	return &n
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
	q := s.queries.WithTx(tx)

	id, err := q.InsertExchange(ctx, db.InsertExchangeParams{
		OrganizerID:  organizerID,
		Name:         name,
		ExchangeDate: date.Format("2006-01-02"),
		State:        string(apitypes.ExchangeStateDraft),
		InviteToken:  token,
		Description:  description,
		BudgetCents:  int32Ptr(req.BudgetCents),
	})
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("insert exchange: %w", err)
	}

	if err := q.InsertExchangeMember(ctx, db.InsertExchangeMemberParams{
		ExchangeID: id,
		UserID:     organizerID,
	}); err != nil {
		return apitypes.Exchange{}, fmt.Errorf("insert organizer membership: %w", err)
	}

	for _, email := range inviteEmails {
		if err := q.InsertExchangeInvite(ctx, db.InsertExchangeInviteParams{
			ExchangeID: id,
			Email:      email,
		}); err != nil {
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
	q := s.queries.WithTx(tx)

	updated, err := q.UpdateExchange(ctx, db.UpdateExchangeParams{
		ID:           exchangeID,
		OrganizerID:  organizerID,
		Name:         name,
		Description:  description,
		ExchangeDate: date.Format("2006-01-02"),
		BudgetCents:  int32Ptr(req.BudgetCents),
	})
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("update exchange: %w", err)
	}
	if updated == 0 {
		return apitypes.Exchange{}, ErrNotFound
	}

	inviteRows, err := q.ListInviteStates(ctx, exchangeID)
	if err != nil {
		return apitypes.Exchange{}, fmt.Errorf("list invites: %w", err)
	}
	existing := make(map[string]db.ListInviteStatesRow, len(inviteRows))
	for _, row := range inviteRows {
		existing[row.Email] = row
	}

	wanted := make(map[string]struct{}, len(inviteEmails))
	var toSend []string
	for _, email := range inviteEmails {
		wanted[email] = struct{}{}
		state, ok := existing[email]
		if !ok {
			if err := q.InsertExchangeInvite(ctx, db.InsertExchangeInviteParams{
				ExchangeID: exchangeID,
				Email:      email,
			}); err != nil {
				return apitypes.Exchange{}, fmt.Errorf("insert invite: %w", err)
			}
			toSend = append(toSend, email)
			continue
		}
		if state.Accepted || !state.Declined {
			continue
		}
		if err := q.ReopenDeclinedInvite(ctx, db.ReopenDeclinedInviteParams{
			ExchangeID: exchangeID,
			Email:      email,
		}); err != nil {
			return apitypes.Exchange{}, fmt.Errorf("reopen invite: %w", err)
		}
		toSend = append(toSend, email)
	}

	for email, state := range existing {
		if _, ok := wanted[email]; ok || state.Accepted || state.Declined {
			continue
		}
		if err := q.DeleteOpenInvite(ctx, db.DeleteOpenInviteParams{
			ExchangeID: exchangeID,
			Email:      email,
		}); err != nil {
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
	deleted, err := s.queries.DeleteExchange(ctx, db.DeleteExchangeParams{
		ID:          exchangeID,
		OrganizerID: userID,
	})
	if err != nil {
		return fmt.Errorf("delete exchange: %w", err)
	}
	if deleted == 0 {
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
