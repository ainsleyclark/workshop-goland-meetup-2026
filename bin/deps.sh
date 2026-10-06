#!/bin/sh
set -eu

log=$(mktemp)
trap 'rm -f "$log"' EXIT

# install runs go install in the background so there is something to watch:
# a first install on a cold module cache takes minutes and is otherwise silent.
# Off a terminal (CI, the student sim) it prints one plain line per tool.
install() {
	name=$1
	pkg=$2
	if [ ! -t 1 ]; then
		echo "  Installing $name"
		go install "$pkg"
		return
	fi
	go install "$pkg" >"$log" 2>&1 &
	pid=$!
	i=0
	while kill -0 "$pid" 2>/dev/null; do
		case $((i % 4)) in
		0) c='|' ;; 1) c='/' ;; 2) c='-' ;; *) c='\' ;;
		esac
		printf '\r  %s Installing %s' "$c" "$name"
		i=$((i + 1))
		sleep 0.1
	done
	if wait "$pid"; then
		printf '\r  \033[32m✓\033[0m %-40s\n' "$name"
	else
		printf '\r  \033[31m✗\033[0m %-40s\n' "$name"
		cat "$log"
		exit 1
	fi
}

printf '\n  Installing the workshop toolchain. The first run can take a few minutes.\n\n'

# Code generation and hot reload.
install sqlc github.com/sqlc-dev/sqlc/cmd/sqlc@latest
install templ github.com/a-h/templ/cmd/templ@v0.3.1020 # In sync with go.mod.
install templier github.com/romshark/templier@latest
install hey github.com/rakyll/hey@latest

# Linting and security checks.
install golangci-lint github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
install gosec github.com/securego/gosec/v2/cmd/gosec@latest
install govulncheck golang.org/x/vuln/cmd/govulncheck@latest

# Claude Code is needed for the Theming block rather than for setup. Its
# installer is per-user and needs no sudo, so it's offered here, but only when
# someone is there to answer: off a terminal (CI, the student sim) nothing is
# downloaded. A failed install never fails setup, which still has upstream to add.
claude_docs=https://code.claude.com/docs/en/setup
if ! command -v claude >/dev/null 2>&1; then
	if [ -t 0 ] && [ -t 1 ]; then
		printf "\n  Claude Code isn't installed. You'll need it for the Theming block.\n"
		printf '  Install it now? [Y/n] '
		read -r answer || answer=n
		case $answer in
		[nN]*) echo "  Skipped. Install it before the Theming block: $claude_docs" ;;
		*)
			if curl -fsSL https://claude.ai/install.sh | bash; then
				command -v claude >/dev/null 2>&1 || [ ! -x "$HOME/.local/bin/claude" ] ||
					echo "  Claude Code is in ~/.local/bin. Add it to your PATH, or open a new terminal."
			else
				echo "  Couldn't install Claude Code just now. Try again later: $claude_docs"
			fi
			;;
		esac
	else
		echo "Claude Code isn't installed yet. You'll need it for the Theming block: $claude_docs"
	fi
fi

# Python is only for make design-search, and installing it means a system
# package manager and often sudo, so name the command for this machine instead.
python_hint() {
	case $(uname -s) in
	Darwin)
		if command -v brew >/dev/null 2>&1; then echo "brew install python"; else echo "xcode-select --install"; fi
		;;
	Linux)
		distro=$(. /etc/os-release 2>/dev/null && echo "$ID ${ID_LIKE:-}")
		case $distro in
		*debian* | *ubuntu*) echo "sudo apt install python3" ;;
		*fedora* | *rhel*) echo "sudo dnf install python3" ;;
		*arch*) echo "sudo pacman -S python" ;;
		*) echo "https://www.python.org/downloads/" ;;
		esac
		;;
	MINGW* | MSYS* | CYGWIN*) echo "winget install Python.Python.3.12" ;;
	*) echo "https://www.python.org/downloads/" ;;
	esac
}
command -v python3 >/dev/null 2>&1 || command -v python >/dev/null 2>&1 ||
	echo "Python 3 isn't installed yet. make design-search needs it: $(python_hint)"
