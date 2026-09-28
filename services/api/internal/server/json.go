package server

import (
	"encoding/json"
	"net/http"

	"github.com/mwritter/giftexchanger/services/api/internal/apitypes"
)

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, apitypes.ErrorResponse{Error: msg})
}
