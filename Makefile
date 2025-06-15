.PHONY: build test lint clean install

# Get version information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'

# Build target
build:
	go build -ldflags "$(LDFLAGS)" -o eolfmt ./cmd/eolfmt

# Run tests
test:
	go test -v ./...

# Static analysis
lint:
	@if command -v revive > /dev/null; then \
		revive -formatter stylish ./...; \
	else \
		echo "revive is not installed. Run 'make install-tools' first."; \
		exit 1; \
	fi

# Install tools
install-tools:
	aqua i -l

# Install binary
install:
	go install -ldflags "$(LDFLAGS)" ./cmd/eolfmt

# Clean up
clean:
	rm -f eolfmt
	go clean

# Run all checks
check: lint test

# Run example
run: build
	./eolfmt .
