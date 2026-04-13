# Makefile for work tool

# Variables
BINARY_COLLECT=collect
BINARY_ANALYZE=analyze
BUILD_DIR=bin
CMD_COLLECT=./cmd/collect
CMD_ANALYZE=./cmd/analyze

.PHONY: all build build-collect build-analyze test test-integration test-all lint fmt vet clean help

all: build test

## build: Build both collect and analyze binaries
build: build-collect build-analyze

## build-collect: Build the collect binary
build-collect: fmt vet
	@echo "Building $(BINARY_COLLECT)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_COLLECT) $(CMD_COLLECT)/main.go
	@chmod +x $(BUILD_DIR)/$(BINARY_COLLECT)

## build-analyze: Build the analyze binary
build-analyze: fmt vet
	@echo "Building $(BINARY_ANALYZE)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_ANALYZE) $(CMD_ANALYZE)/main.go
	@chmod +x $(BUILD_DIR)/$(BINARY_ANALYZE)

## test: Run unit tests
test: 
	@echo "Running unit tests..."
	@go test -short -count=1 -cover ./...

## test-integration: Run integration tests (requires credentials)
test-integration:
	@echo "Running integration tests..."
	@go test -v -count=1 ./...

## test-all: Run all tests (unit + integration)
test-all:
	@echo "Running all tests..."
	@go test -v -count=1 ./...

## lint: Run basic linting (vet and fmt check)
lint: vet fmt-check

## fmt: Format all Go files
fmt:
	@echo "Formatting Go files..."
	@go fmt ./...

## fmt-check: Check if Go files are formatted
fmt-check:
	@echo "Checking Go files formatting..."
	@test -z $$(gofmt -l .) || (echo "Files not formatted. Run 'make fmt'"; exit 1)

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(CMD_COLLECT)/$(BINARY_COLLECT)

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' Makefile | column -t -s ':' |  sed -e 's/^/  /'
