package api

import "net/http"

// HandlerOptions wires dependencies for the HTTP API.
type HandlerOptions struct {
	ServiceName string
	Version     VersionService
}

// NewHandler builds the HTTP handler tree for the backend service.
func NewHandler(opts HandlerOptions) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/healthz", HealthHandler{ServiceName: opts.ServiceName})
	if opts.Version != nil {
		versionHandler := VersionHandler{Service: opts.Version}
		mux.Handle("/api/v1/version/latest", versionHandler)
		mux.Handle("/api/v1/version/", versionHandler)
	}
	return mux
}
