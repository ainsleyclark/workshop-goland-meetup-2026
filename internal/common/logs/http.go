package logs

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestLogger logs each request using the supplied logger.
//
// Place it outside recovery middleware to record recovered responses.
// Unrecovered panics continue to the server and are logged as
// interrupted requests, not successful responses.
func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			start := time.Now()
			wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			completed := false

			defer func() {
				status := wrapped.Status()
				if status == 0 && completed {
					status = http.StatusOK
				}

				level, message := slog.LevelInfo, "HTTP request"
				switch {
				case !completed:
					level, message = slog.LevelError, "HTTP request interrupted"
				case status >= 500:
					level = slog.LevelError
				case status >= 400:
					level = slog.LevelWarn
				}

				log.LogAttrs(ctx, level, message,
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", status),
					slog.Int("bytes", wrapped.BytesWritten()),
					slog.Duration("duration", time.Since(start)),
				)
			}()

			next.ServeHTTP(wrapped, r)

			completed = true
		})
	}
}
