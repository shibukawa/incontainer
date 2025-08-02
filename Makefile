.PHONY: build test clean docker-build docker-test test-all

# Build the CLI tool
build:
	go build -o iscontainer ./cmd/iscontainer

# Run Go tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f iscontainer
	docker rmi incontainer-test 2>/dev/null || true

# Build Docker image for testing
docker-build:
	docker build -t incontainer-test .

# Test with Docker
docker-test: docker-build
	@echo "=== Testing with Docker ==="
	docker run --rm incontainer-test

# Test with Podman (if available)
podman-test:
	@echo "=== Building with Podman ==="
	podman build -t incontainer-test .
	@echo "=== Testing with Podman ==="
	podman run --rm incontainer-test

# Test locally
local-test: build
	@echo "=== Testing on host system ==="
	./iscontainer -v

# Run all available tests
test-all: test local-test docker-test
	@echo "All tests completed!"

# Install the CLI tool
install: build
	cp iscontainer /usr/local/bin/

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the CLI tool"
	@echo "  test         - Run Go tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-test  - Test with Docker"
	@echo "  podman-test  - Test with Podman"
	@echo "  local-test   - Test on host system"
	@echo "  test-all     - Run all tests"
	@echo "  install      - Install CLI tool to /usr/local/bin"
	@echo "  help         - Show this help"