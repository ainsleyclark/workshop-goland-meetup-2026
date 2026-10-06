#!/bin/sh
set -eu

# Code generation and hot reload.
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/a-h/templ/cmd/templ@v0.3.1020 # In sync with go.mod.
go install github.com/romshark/templier@latest
go install github.com/rakyll/hey@latest

# Linting and security checks.
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest

# Needed for the Theming block rather than for setup, so only warn.
command -v claude >/dev/null 2>&1 ||
	echo "Claude Code isn't installed yet. You'll need it for the Theming block:" \
		"https://code.claude.com/docs/en/setup"

command -v python3 >/dev/null 2>&1 || command -v python >/dev/null 2>&1 ||
	echo "Python 3 isn't installed yet. make design-search needs it."
