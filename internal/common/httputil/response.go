package httputil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// WriteJSON writes v as JSON without HTML escaping or a trailing newline.
// If encoding fails, it writes an HTTP 500 response instead.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", applicationJSON)
	w.WriteHeader(status)
	_, _ = w.Write(bytes.TrimSuffix(buf.Bytes(), []byte("\n")))
}

// WriteRawJSON writes pre-encoded JSON unchanged with content sniffing disabled.
// The caller is responsible for providing valid JSON.
func WriteRawJSON(w http.ResponseWriter, raw string) {
	w.Header().Set("Content-Type", applicationJSON)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write([]byte(raw)) // #nosec G705 -- served as JSON with nosniff, never HTML.
}
