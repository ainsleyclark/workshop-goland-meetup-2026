# The workshop repository this copy was forked from, as owner/repo. make publish
# in the maintainer repository rewrites this line for each event.
UPSTREAM := ainsleyclark/workshop-goland-meetup-2026

# Fails when origin is the workshop repository itself rather than a fork of it,
# so a push never 403s after an hour of work.
define require_fork
	@case "$$(git remote get-url origin 2>/dev/null)" in *$(UPSTREAM)*) \
		echo "This is a clone of the workshop repository, not a fork of it."; \
		echo "Fork https://github.com/$(UPSTREAM) on GitHub, clone your fork, and run this there."; exit 1;; esac
endef

setup: # Install the toolchain and point upstream at the workshop repository; run once after cloning your fork
	$(require_fork)
	sh -e bin/deps.sh
	@git remote get-url upstream >/dev/null 2>&1 || git remote add upstream https://github.com/$(UPSTREAM).git
	@git fetch -q upstream || echo "Couldn't reach $(UPSTREAM) just now; make update fetches it later."
	@printf '\n  ==============================================\n'
	@printf '    Welcome to the workshop!\n'
	@printf '    Full-Stack Go with sqlc and templ\n'
	@printf '  ==============================================\n\n'
	@printf '  You are ready. Try make run, or make web for the site on :8080.\n'
	@printf '  Have a look through ./architectures before you start.\n\n'
.PHONY: setup

update: # Pull the instructor's latest changes from the workshop repository into your fork
	git pull --rebase upstream main
.PHONY: update

run: # Run the program
	go run .
.PHONY: run

run-local: # Run the program against the local mock APIs (make mock), for when they're down
	go run . --local
.PHONY: run-local

web: # Serve the front-end
	go run main.go web
.PHONY: web

web-hot: # Serve hot reload with templ
	CLICOLOR_FORCE=1 templier --config ./templier.yml
.PHONY: web-hot

key: # Save the workshop's Anthropic key for Claude Code, then run `claude` here and type /design-site
	@py=$$(command -v python3 || command -v python) || { echo "make key needs Python 3"; exit 1; }; \
	printf 'Paste the workshop key from Slack, or press Enter to use your own Claude login: '; \
	read -r key; \
	[ -n "$$key" ] || { echo "No key saved; Claude Code will use your own login."; exit 0; }; \
	KEY="$$key" "$$py" -c 'import json, os, pathlib; p = pathlib.Path(".claude/settings.local.json"); s = json.loads(p.read_text()) if p.exists() else {}; s.setdefault("env", {})["ANTHROPIC_API_KEY"] = os.environ["KEY"].strip().strip("\x27\x22"); p.parent.mkdir(exist_ok=True); p.write_text(json.dumps(s, indent=2) + "\n")' && \
	echo "Saved to .claude/settings.local.json. Start Claude Code here with: claude, then type /design-site"
.PHONY: key

design-search: # Search the ui-ux-pro-max design data, e.g. make design-search Q="warm and earthy"
	@[ -n "$(Q)" ] || { echo 'Usage: make design-search Q="what you want" [ARGS="--domain color"]'; exit 1; }
	@py=$$(command -v python3 || command -v python) || { echo "design-search needs Python 3"; exit 1; }; \
	PYTHONDONTWRITEBYTECODE=1 "$$py" .claude/skills/ui-ux-pro-max/scripts/search.py "$(Q)" $(or $(ARGS),--design-system)
.PHONY: design-search

mock: # Serve the mock GBIF and Open-Meteo APIs on :8082, for when they're down
	go run ./cmd/mockapi
.PHONY: mock

mock-scrape: # Refresh the mock's GBIF data from the live API
	go run ./cmd/mockapi scrape
.PHONY: mock-scrape

submit: # Commit everything and push to your fork; opens a draft PR to the workshop repository if gh is installed
	$(require_fork)
	git add -A .
	git commit -m "Workshop progress" || true
	git push -u origin HEAD
	@command -v gh >/dev/null 2>&1 || exit 0; \
	gh pr view --repo $(UPSTREAM) >/dev/null 2>&1 && exit 0; \
	gh pr create --draft --repo $(UPSTREAM) --head "$$(gh api user --jq .login):$$(git branch --show-current)" \
		--title "$$(git config user.name)'s sightings site" \
		--body "Workshop progress, for the instructor to look at. Not for merging." >/dev/null 2>&1 \
		&& echo "Opened a draft pull request on $(UPSTREAM) so the instructor can see your site." || true
.PHONY: submit

format: # Run gofmt and templ fmt
	go fmt ./...
	go tool templ fmt .
.PHONY: format

test: # Run all tests
	go clean -testcache && go test ./... -coverprofile=coverage.out -covermode=atomic
.PHONY: test

test-race: # Run all tests with race
	go clean -testcache && go test -race ./... -coverprofile=coverage.out -covermode=atomic
.PHONY: test-race

load: # Load test the home page for 20 seconds while make web runs
	@command -v hey >/dev/null 2>&1 || go install github.com/rakyll/hey@latest
	hey -z 10s -c 10 http://localhost:8080/
.PHONY: load

load-pprof: # Load test the home page for 20 seconds while make web runs, saving CPU and allocation profiles here
	@command -v hey >/dev/null 2>&1 || go install github.com/rakyll/hey@latest
	@curl -s -o ./cpu.pprof "http://localhost:6060/debug/pprof/profile?seconds=10" & \
	hey -z 20s -c 10 http://localhost:8080/; \
	wait; \
	curl -s -o ./allocs.pprof http://localhost:6060/debug/pprof/allocs; \
	echo; \
	echo "Profiles saved, open one with:"; \
	echo "  go tool pprof -http=: ./cpu.pprof"; \
	echo "  go tool pprof -http=: -sample_index=alloc_space ./allocs.pprof"
.PHONY: load-pprof

lint: # Run linter
	golangci-lint run ./... --fix --config=.golangci.yaml
.PHONY: lint

clean:
	go fmt ./...
	go fix ./...
.PHONY: clean

gen: # Runs all //go:generate
	go generate ./...
.PHONY: gen

templ: # Runs templ generate
	go tool templ generate ./...
.PHONY: templ

# sqlc on PATH, or where bin/deps.sh installed it for anyone without ~/go/bin on their PATH.
SQLC := $(or $(shell command -v sqlc 2>/dev/null),$(shell go env GOPATH)/bin/sqlc)

sqlc: # Regenerate sqlc output, then build, so a query change that breaks a store shows up now
	$(SQLC) generate
	go build ./...
.PHONY: sqlc

sec: # Run gosec security scan (matches CI)
	@command -v gosec >/dev/null 2>&1 || go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -exclude-generated ./...
.PHONY: sec

vuln: # Run govulncheck (matches CI)
	@command -v govulncheck >/dev/null 2>&1 || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...
.PHONY: vuln

kill-ports: # Kill anything listening on the ports this app uses (8080 web, 8081 web-hot, 8082 mock, 8083 themes)
	@for port in 8080 8081 8082 8083; do \
		pid=$$(lsof -ti tcp:$$port); \
		if [ -n "$$pid" ]; then \
			echo "Killing $$pid on :$$port"; \
			kill -9 $$pid; \
		fi; \
	done
.PHONY: kill-ports

all: # Make format, lint and test
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) templ
	$(MAKE) test-race
	$(MAKE) vuln
	$(MAKE) sec
.PHONY: all

todo: # Show to-do items per file
	$(Q) grep \
		--exclude=Makefile.util \
		--exclude=Makefile \
		--exclude=TODO.md \
		--exclude-dir=vendor \
		--exclude-dir=.vercel \
		--exclude-dir=.gen \
		--exclude-dir=.idea \
		--exclude-dir=public \
		--exclude-dir=node_modules \
		--exclude-dir=archetypes \
		--exclude-dir=.git \
		--text \
		--color \
		-nRo \
		-E '\S*[^\.]TODO.*' \
		.
.PHONY: todo

help: # Display this help
	$(Q) awk 'BEGIN {FS = ":.*#"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?#/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help
