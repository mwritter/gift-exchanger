package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mwritter/giftexchanger/services/api/internal/apitypes"
	"github.com/mwritter/giftexchanger/services/api/internal/exchanges"
)

func (s *Server) getExchanges(w http.ResponseWriter, r *http.Request) {
	// return all exchanges the current user is a part of or is the organizer for
	// for now if you organized the exchange you are apart of the exchange
	user := UserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := s.exchanges.GetExchanges(r.Context(), user.ID)
	if err != nil {
		log.Printf("list exchanges: %v", err)
		writeError(w, http.StatusInternalServerError, "could not list exchanges")
		return
	}

	writeJSON(w, http.StatusOK, apitypes.ListExchangesResponse{
		Exchanges: result,
	})
}

func (s *Server) getExchange(w http.ResponseWriter, r *http.Request) {
	// get exchange by id - must be organizer or member
	user := UserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := s.exchanges.GetExchangeById(r.Context(), user.ID, chi.URLParam(r, "exchangeID"))
	if errors.Is(err, exchanges.ErrNotFound) {
		writeError(w, http.StatusNotFound, "exchange not found")
		return
	}
	if err != nil {
		log.Printf("get exchange: %v", err)
		writeError(w, http.StatusInternalServerError, "could not get exchange")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) getExchangeInvites(w http.ResponseWriter, r *http.Request) {
	// list the open invites for an exchange - only the organizer can do this
	user := UserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := s.exchanges.GetExchangeInvites(r.Context(), user.ID, chi.URLParam(r, "exchangeID"))
	if errors.Is(err, exchanges.ErrNotFound) {
		writeError(w, http.StatusNotFound, "exchange not found")
		return
	}
	if errors.Is(err, exchanges.ErrForbidden) {
		writeError(w, http.StatusForbidden, "only the organizer can see invites")
		return
	}
	if err != nil {
		log.Printf("list exchange invites: %v", err)
		writeError(w, http.StatusInternalServerError, "could not list invites")
		return
	}

	writeJSON(w, http.StatusOK, apitypes.ListExchangeInvitesResponse{
		InviteEmails: result,
	})
}

func (s *Server) createExchange(w http.ResponseWriter, r *http.Request) {
	// create exchange, add current user as the organizer and as a member
	var req apitypes.CreateExchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := s.exchanges.CreateExchange(r.Context(), user.ID, user.Email, req)
	if errors.Is(err, exchanges.ErrInvalidName) {
		writeError(w, http.StatusBadRequest, "Give your exchange a name")
		return
	}
	if errors.Is(err, exchanges.ErrInvalidDate) {
		writeError(w, http.StatusBadRequest, "Pick an exchange date")
		return
	}
	if errors.Is(err, exchanges.ErrInvalidBudget) {
		writeError(w, http.StatusBadRequest, "Budget must be zero or more")
		return
	}
	if errors.Is(err, exchanges.ErrInvalidEmail) {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if err != nil {
		log.Printf("create exchange: %v", err)
		writeError(w, http.StatusInternalServerError, "could not create exchange")
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (s *Server) updateExchange(w http.ResponseWriter, r *http.Request) {
	// update exchange by id - only the orgainzer can do this
	var req apitypes.UpdateExchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := s.exchanges.UpdateExchange(r.Context(), chi.URLParam(r, "exchangeID"), user.ID, user.Email, req)
	if errors.Is(err, exchanges.ErrNotFound) {
		writeError(w, http.StatusNotFound, "exchange not found")
		return
	}
	if errors.Is(err, exchanges.ErrInvalidName) {
		writeError(w, http.StatusBadRequest, "Give your exchange a name")
		return
	}
	if errors.Is(err, exchanges.ErrInvalidDate) {
		writeError(w, http.StatusBadRequest, "Pick an exchange date")
		return
	}
	if errors.Is(err, exchanges.ErrInvalidBudget) {
		writeError(w, http.StatusBadRequest, "Budget must be zero or more")
		return
	}
	if errors.Is(err, exchanges.ErrInvalidEmail) {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if err != nil {
		log.Printf("update exchange: %v", err)
		writeError(w, http.StatusInternalServerError, "could not update exchange")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) deleteExchange(w http.ResponseWriter, r *http.Request) {
	// delete exchange by id - only the orgainzer can do this
	user := UserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := s.exchanges.DeleteExchange(r.Context(), chi.URLParam(r, "exchangeID"), user.ID)
	if errors.Is(err, exchanges.ErrNotFound) {
		writeError(w, http.StatusNotFound, "exchange not found")
		return
	}
	if err != nil {
		log.Printf("delete exchange: %v", err)
		writeError(w, http.StatusInternalServerError, "could not delete exchange")
		return
	}

	writeJSON(w, http.StatusOK, apitypes.StatusResponse{Status: "ok"})
}
func (s *Server) acceptExchangeInvite(w http.ResponseWriter, r *http.Request) {
	// accept invite
}
func (s *Server) declineExchangeInvite(w http.ResponseWriter, r *http.Request) {
	// decline invite
}
