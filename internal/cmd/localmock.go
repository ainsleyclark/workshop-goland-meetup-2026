package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"workshop/internal/common/printer"
)

// localMockAddr is the address the local mock API server (cmd/mockapi)
// listens on, matching localGBIFBaseURL/localOpenMeteoBaseURL.
const localMockAddr = "localhost:8082"

// localMockStartTimeout bounds how long ensureLocalMock waits for a
// freshly started mock server to come up.
const localMockStartTimeout = 20 * time.Second

// ensureLocalMock starts the local mock API server in the background if
// nothing is already listening on localMockAddr.
func ensureLocalMock(ctx context.Context, out *printer.Console) error {
	if dialAddr(localMockAddr, time.Second) {
		return nil
	}

	out.Info("Starting the local mock API in the background...")

	logPath := filepath.Join(os.TempDir(), "workshop-mockapi.log")
	logFile, err := os.Create(logPath) // #nosec G304 -- fixed filename under the OS temp dir, not user input.
	if err != nil {
		return fmt.Errorf("creating mock API log file: %w", err)
	}

	cmd := exec.Command("go", "run", "./cmd/mockapi")
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting mock API: %w", err)
	}

	deadline := time.After(localMockStartTimeout)
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		if dialAddr(localMockAddr, time.Second) {
			out.Infof("Local mock API is up (logs: %s).\n", logPath)
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline:
			return fmt.Errorf("mock API didn't come up within %s, check %s", localMockStartTimeout, logPath)
		case <-ticker.C:
		}
	}
}

// dialAddr reports whether something is already listening at addr.
func dialAddr(addr string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
