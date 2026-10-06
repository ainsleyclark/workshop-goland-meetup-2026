package gbif

import (
	"net/http"
	"time"
	"workshop/internal/common/httputil"
)

// writeError writes an error in the shape GBIF uses.
func writeError(w http.ResponseWriter, status int, message string) {
	httputil.WriteJSON(w, status, struct {
		Timestamp string `json:"timestamp"`
		Status    int    `json:"status"`
		Error     string `json:"error"`
		Message   string `json:"message"`
	}{
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000+00:00"),
		Status:    status,
		Message:   message,
	})
}
