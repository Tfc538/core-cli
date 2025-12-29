package api

import (
	"net/http"
	"strings"
)

const versionPrefix = "/api/v1/version/"

// VersionHandler serves version metadata endpoints.
type VersionHandler struct {
	Service VersionService
}

func (h VersionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	path := r.URL.Path
	if path == "/api/v1/version/latest" {
		info, err := h.Service.Latest(r.Context())
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to load latest version")
			return
		}
		WriteJSON(w, http.StatusOK, Response{Status: "ok", Data: info})
		return
	}

	if !strings.HasPrefix(path, versionPrefix) {
		WriteError(w, http.StatusNotFound, "not found")
		return
	}

	version := strings.TrimPrefix(path, versionPrefix)
	if version == "" {
		WriteError(w, http.StatusNotFound, "not found")
		return
	}

	info, ok, err := h.Service.Get(r.Context(), version)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to load version")
		return
	}
	if !ok {
		WriteError(w, http.StatusNotFound, "version not found")
		return
	}

	WriteJSON(w, http.StatusOK, Response{Status: "ok", Data: info})
}
