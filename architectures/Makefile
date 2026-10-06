# Each architecture is its own Go module, so every target walks them in turn.
modules := $(patsubst %/go.mod,%,$(wildcard */go.mod))

# Binaries go to a scratch directory; nothing generated should land in a folder.
out := $(shell mktemp -d)

check: # Tidy, format, and confirm every architecture compiles
	@for m in $(modules); do \
		echo "==> $$m"; \
		(cd $$m && go mod tidy && go fmt ./... && go build -o $(out)/ ./... && go vet ./... && go fix ./...) || exit 1; \
	done
	@echo "==> All architectures compile & checked"
.PHONY: check

help: # Display this help
	$(Q) awk 'BEGIN {FS = ":.*#"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?#/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help
