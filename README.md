# web-analyzer

Web-application for analyzing Webpages built in Golang.

## Objective

The objective is to build a web application that does an analysis of a
web-page/URL. The application shows a form with a text field in which users can type in the URL of the webpage to be analyzed. Additionally to the form, it contains a button to send a request to the server.

More details about the requirements can be found in the [requirement document](docs/requirement.md).

## Structure

The project is structured as follows:

```tree
.
├── build
│   └── package
│       └── Dockerfile
├── cmd
│   └── server
│       └── main.go
├── internal
│   └── handlers
│       └── handlers.go
├── docs
│   └── requirement.md
├── deploy
│   ├── base
│   │   ├── configmap.yaml
│   │   ├── deployment.yaml
│   │   ├── kustomization.yaml
│   │   ├── namespace.yaml
│   │   └── service.yaml
│   └── overlays
│       ├── dev
│       ├── staging
│       └── prod
├── scripts
│   └── k8s-local.sh
├── go.mod
├── go.sum
├── Makefile
└── README.md

```

For more details, see also the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## Build

run `make build` to build the application.

## Run

run `make run` to start the application. By default, the server will start on port 8080. You can access it by navigating to `http://localhost:8080` in your web browser.

## Test

run `make test` to execute the unit tests for the application.

## Deploy

## Docker

run `make docker-build` to build the Docker image for the application.
run `make docker-run` to run the Docker container for the application.
run `make docker-stop` to stop the Docker container for the application.

## Kubernetes Deployment

The application includes Kubernetes manifests using Kustomize for different environments (dev, staging, prod).

### Prerequisites

- Docker Desktop with Kubernetes enabled (for local development)
- `kubectl` CLI installed
- `kustomize` (optional, kubectl has built-in kustomize support)

### Local Kubernetes Deployment

For local development using Docker Desktop Kubernetes:

```bash
# Deploy to local Kubernetes
make k8s-local

# Or run the script directly
./scripts/k8s-local.sh
```

### Deploy to Different Environments

```bash
# Deploy to dev environment
make k8s-deploy-dev

# Deploy to staging environment
make k8s-deploy-staging

# Deploy to production environment
make k8s-deploy-prod
```

### Other Kubernetes Commands

```bash
# View deployment status
make k8s-status ENV=dev

# Delete deployment
make k8s-delete ENV=dev

# View logs
make k8s-logs ENV=dev
```

### Directory Structure

```
deploy/
├── base/                    # Base Kubernetes manifests
│   ├── configmap.yaml
│   ├── deployment.yaml
│   ├── kustomization.yaml
│   ├── namespace.yaml
│   └── service.yaml
└── overlays/               # Environment-specific overlays
    ├── dev/
    ├── staging/
    └── prod/
```

## CI/CD

### Continuous Integration (GitHub Actions)

The project includes comprehensive CI workflows:

| Workflow | Trigger | Description |
|----------|---------|-------------|
| **CI** | Push/PR to main, develop | Lint, test, build, security scan |
| **Release** | Tag push (v*) | Create GitHub release with binaries and Docker images |
| **CodeQL** | Push/PR to main, weekly | Security and code quality analysis |

#### CI Pipeline Stages

1. **Lint** - golangci-lint, go vet, format check
2. **Test** - Unit tests with race detection and coverage
3. **Build** - Binary compilation for Linux AMD64
4. **Docker** - Build and test Docker image
5. **Security** - Gosec and govulncheck scanning
6. **Dependency Review** - PR dependency analysis

### Automated Dependency Updates

Dependabot is configured for:

- Go modules (weekly)
- GitHub Actions (weekly)
- Docker base images (weekly)

## Improvements

- [ ]  Continuous Deployment (CD) pipeline setup (e.g., Deploy to Azure AKS or AWS EKS using GitHub Actions and Terraform)
- [ ]  Add metrics, tracing, and logging (e.g., Grafana Stack and OpenTelemetry)
- [ ]  AuthN/AuthZ
