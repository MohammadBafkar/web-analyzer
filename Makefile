# Variables
APP_NAME := web-analyzer
BINARY_DIR := bin
BINARY := $(BINARY_DIR)/$(APP_NAME)
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := latest
GO := go
GOFLAGS := -v

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: all build run test test-coverage clean deps fmt lint vet docker-build docker-run docker-stop help

## all: Build and test the application
all: clean deps build test

## build: Build the application binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BINARY)
	$(GO) build $(GOFLAGS) -o $(BINARY) ./...

## run: Run the application
run:
	@echo "Running $(APP_NAME)..."
	$(GO) run $(GOFLAGS) ./...

## test: Run all tests
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

## docker-stop: Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE):$(DOCKER_TAG)) 2>/dev/null || true

## install: Install the binary
install: build
	@echo "Installing $(APP_NAME)..."
	$(GO) install ./...

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'