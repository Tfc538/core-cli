package api

import (
	"encoding/json"
	"net/http"
)

// Response is the standard JSON envelope for the backend API.
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// WriteJSON writes a JSON response with the provided status code.
func WriteJSON(w http.ResponseWriter, status int, payload Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteError writes a standardized error response.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, Response{Status: "error", Error: message})
}
