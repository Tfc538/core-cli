package api

import "net/http"

// HandlerOptions wires dependencies for the HTTP API.
type HandlerOptions struct {
	ServiceName string
	Version     VersionService
}

// NewHandler builds the HTTP handler tree for the backend service.
func NewHandler(opts HandlerOptions) http.Handler {
	_ = opts
	mux := http.NewServeMux()
	return mux
}
