package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standard JSON envelope sent back on every request.
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// SendJSON is a helper function that encodes the response.
func SendJSON(w http.ResponseWriter, statusCode int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
