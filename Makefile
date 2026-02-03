# Variables
APP_NAME := web-analyzer
BINARY_DIR := bin
BINARY := $(BINARY_DIR)/$(APP_NAME)
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := latest
GO := go
GOFLAGS := -v

# Docker variables
DOCKERFILE := build/package/Dockerfile

# Kubernetes variables
K8S_NAMESPACE_DEV := web-analyzer-dev
K8S_NAMESPACE_STAGING := web-analyzer-staging
K8S_NAMESPACE_PROD := web-analyzer-prod
ENV ?= dev

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: all build run test test-coverage clean deps fmt lint vet docker-build docker-run docker-stop help
.PHONY: k8s-local k8s-deploy-dev k8s-deploy-staging k8s-deploy-prod k8s-delete k8s-status k8s-logs

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
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f $(DOCKERFILE) .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

## docker-stop: Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE):$(DOCKER_TAG)) 2>/dev/null || true

## k8s-local: Deploy to local Kubernetes (Docker Desktop)
k8s-local:
	@echo "Deploying to local Kubernetes..."
	@chmod +x ./scripts/k8s-local.sh
	./scripts/k8s-local.sh

## k8s-deploy-dev: Deploy to dev environment
k8s-deploy-dev: docker-build
	@echo "Deploying to dev environment..."
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):dev
	kubectl apply -k deploy/overlays/dev
	kubectl rollout status deployment/dev-web-analyzer -n $(K8S_NAMESPACE_DEV) --timeout=120s

## k8s-deploy-staging: Deploy to staging environment
k8s-deploy-staging: docker-build
	@echo "Deploying to staging environment..."
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):staging
	kubectl apply -k deploy/overlays/staging
	kubectl rollout status deployment/staging-web-analyzer -n $(K8S_NAMESPACE_STAGING) --timeout=120s

## k8s-deploy-prod: Deploy to production environment
k8s-deploy-prod: docker-build
	@echo "Deploying to production environment..."
	kubectl apply -k deploy/overlays/prod
	kubectl rollout status deployment/prod-web-analyzer -n $(K8S_NAMESPACE_PROD) --timeout=120s

## k8s-delete: Delete Kubernetes deployment (use ENV=dev|staging|prod)
k8s-delete:
	@echo "Deleting $(ENV) deployment..."
	kubectl delete -k deploy/overlays/$(ENV) --ignore-not-found=true

## k8s-status: Show deployment status (use ENV=dev|staging|prod)
k8s-status:
	@echo "Deployment status for $(ENV):"
	@kubectl get all -n web-analyzer-$(ENV)

## k8s-logs: View application logs (use ENV=dev|staging|prod)
k8s-logs:
	@kubectl logs -f -l app=web-analyzer -n web-analyzer-$(ENV)

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