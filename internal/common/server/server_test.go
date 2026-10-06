package server_test

import (
	"context"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
	"workshop/internal/common/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, handler http.Handler) (string, func() error) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "binding an ephemeral port")

	ctx, cancel := context.WithCancel(t.Context())
	errCh := make(chan error, 1)

	srv := server.New(0, handler)
	go func() {
		errCh <- srv.Serve(ctx, listener)
	}()

	stop := sync.OnceValue(func() error {
		cancel()
		return <-errCh
	})
	t.Cleanup(func() { _ = stop() })

	return "http://" + listener.Addr().String(), stop
}

func TestServer_Serve(t *testing.T) {
	t.Parallel()

	t.Run("Serves the handler", func(t *testing.T) {
		t.Parallel()

		url, _ := setup(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}))

		res, err := http.Get(url)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
	})

	t.Run("Cancelling the context is a clean shutdown", func(t *testing.T) {
		t.Parallel()

		_, stop := setup(t, http.NotFoundHandler())

		assert.NoError(t, stop())
	})

	t.Run("Stops serving once shut down", func(t *testing.T) {
		t.Parallel()

		url, stop := setup(t, http.NotFoundHandler())
		require.NoError(t, stop())

		_, err := http.Get(url)
		assert.Error(t, err)
	})
}

func TestServer_Start(t *testing.T) {
	t.Parallel()

	t.Run("Fails when the port is already bound", func(t *testing.T) {
		t.Parallel()

		// Bind every interface, as Start does: holding only 127.0.0.1
		// does not stop a later bind of ":port" on macOS.
		listener, err := net.Listen("tcp", ":0")
		require.NoError(t, err)
		defer listener.Close()

		port := listener.Addr().(*net.TCPAddr).Port
		err = server.New(port, http.NotFoundHandler()).Start(t.Context())

		assert.ErrorContains(t, err, "listening on")
	})
}

func TestServer_ShutdownClosesActiveConnectionsAfterDeadline(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	finished := make(chan struct{})
	url, stop := setup(t, http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
		close(finished)
	}))
	client := &http.Client{Timeout: 10 * time.Second}
	requestDone := make(chan error, 1)
	go func() {
		res, err := client.Get(url)
		if res != nil {
			res.Body.Close()
		}
		requestDone <- err
	}()

	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("request did not reach the handler")
	}
	require.NoError(t, stop())
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("shutdown left the active connection open")
	}
	require.Error(t, <-requestDone)
}
