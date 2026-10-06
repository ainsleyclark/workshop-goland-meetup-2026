package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Server serves an http.Handler for as long as a context lives.
type Server struct {
	srv *http.Server
}

// New builds a Server that serves handler on the given port.
func New(port int, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           handler,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

const (
	// readHeaderTimeout bounds how long a client may take to send its
	// headers, which is the cheapest guard against a slowloris.
	readHeaderTimeout = 5 * time.Second

	// shutdownTimeout caps how long in-flight requests get to finish once
	// a shutdown begins. It is a ceiling, not a wait: draining returns as
	// soon as the last request does.
	shutdownTimeout = 2 * time.Second
)

// Start binds the configured port and serves until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("listening on %s: %w", s.srv.Addr, err)
	}

	return s.Serve(ctx, ln)
}

// Serve is Start with a listener the caller already holds, which lets
// tests bind port zero and discover the address that was chosen.
//
// It returns nil once the server has stopped cleanly, so a cancelled
// context is a successful shutdown rather than an error.
func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.srv.Serve(ln)
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serving http: %w", err)
	case <-ctx.Done():
		// The parent is already cancelled, so draining needs its own deadline.
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
		defer cancel()

		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			// Shutdown leaves connections open when its deadline expires.
			// Force them closed so an interrupted command can still exit.
			if closeErr := s.srv.Close(); closeErr != nil {
				return fmt.Errorf("closing http server: %w", closeErr)
			}
			if !errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("shutting down http server: %w", err)
			}
		}

		return nil
	}
}
