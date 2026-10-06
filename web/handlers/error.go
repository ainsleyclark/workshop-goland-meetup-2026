package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"workshop/web/views/pages"
)

// renderError writes code as the response status and renders the error page
// with message. It logs if rendering itself fails.
func renderError(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, code int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	if err := pages.Error(message, code).Render(ctx, w); err != nil {
		logger.ErrorContext(ctx, "handlers: rendering error page: "+err.Error())
	}
}
