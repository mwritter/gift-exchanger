// Package apitypes holds the JSON request and response shapes of the HTTP API.
// TypeScript types in packages/types are generated from this package with
// tygo (`pnpm types:generate`); edit here, not in the generated file.
package apitypes

type ErrorResponse struct {
	Error string `json:"error"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type MagicLinkRequest struct {
	Email string `json:"email"`
}

type VerifyLoginCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ExchangeState string

const (
	ExchangeStateDraft     ExchangeState = "draft"
	ExchangeStateOpen      ExchangeState = "open"
	ExchangeStateActive    ExchangeState = "active"
	ExchangeStateCompleted ExchangeState = "completed"
)

type Exchange struct {
	ID          string `json:"id"`
	OrganizerID string `json:"organizerId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Calendar date, YYYY-MM-DD.
	ExchangeDate string        `json:"exchangeDate"`
	BudgetCents  *int          `json:"budgetCents" tstype:"number | null"`
	State        ExchangeState `json:"state"`
	InviteEmails []string      `json:"inviteEmails"`
	MemberCount  int           `json:"memberCount"`
}

type CreateExchangeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Calendar date, YYYY-MM-DD.
	ExchangeDate string   `json:"exchangeDate"`
	BudgetCents  *int     `json:"budgetCents" tstype:"number | null"`
	InviteEmails []string `json:"inviteEmails"`
}

type UpdateExchangeRequest = CreateExchangeRequest

type ListExchangesResponse struct {
	Exchanges []Exchange `json:"exchanges"`
}

// Invite is a pending invite for the logged-in user's email.
type Invite struct {
	ExchangeID   string `json:"exchangeId"`
	ExchangeName string `json:"exchangeName"`
	// Calendar date, YYYY-MM-DD.
	ExchangeDate   string `json:"exchangeDate"`
	OrganizerName  string `json:"organizerName"`
	OrganizerEmail string `json:"organizerEmail"`
	InvitedAt      string `json:"invitedAt"`
}

type ListInvitesResponse struct {
	Invites []Invite `json:"invites"`
}
