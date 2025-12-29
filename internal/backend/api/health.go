package api

import "net/http"

type HealthHandler struct {
	ServiceName string
}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	WriteJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data: map[string]string{
			"service": h.ServiceName,
			"status":  "ok",
		},
	})
}
