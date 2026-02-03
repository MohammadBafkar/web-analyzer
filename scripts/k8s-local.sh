#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo -e "${GREEN}🚀 Deploying Web Analyzer to local Kubernetes (Docker Desktop)${NC}"

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}❌ kubectl is not installed. Please install it first.${NC}"
    exit 1
fi

# Check if Docker Desktop Kubernetes is running
if ! kubectl config get-contexts docker-desktop &> /dev/null; then
    echo -e "${RED}❌ docker-desktop context not found. Please enable Kubernetes in Docker Desktop.${NC}"
    exit 1
fi

if ! kubectl --context docker-desktop cluster-info &> /dev/null; then
    echo -e "${RED}❌ Docker Desktop Kubernetes is not running. Please start Docker Desktop with Kubernetes enabled.${NC}"
    exit 1
fi

# Switch to docker-desktop context
echo -e "${YELLOW}📌 Switching to docker-desktop context...${NC}"
kubectl config use-context docker-desktop

# Build the Docker image
echo -e "${YELLOW}🔨 Building Docker image...${NC}"
docker build -t web-analyzer:local -f "${PROJECT_ROOT}/build/package/Dockerfile" "${PROJECT_ROOT}"
docker tag web-analyzer:local web-analyzer:dev

# Create namespace if it doesn't exist
echo -e "${YELLOW}📦 Creating namespace...${NC}"
kubectl create namespace web-analyzer-dev --dry-run=client -o yaml | kubectl apply -f -

# Apply Kustomize overlay for dev
echo -e "${YELLOW}📋 Applying Kubernetes manifests...${NC}"
kubectl apply -k "${PROJECT_ROOT}/deploy/overlays/dev"

# Wait for deployment
echo -e "${YELLOW}⏳ Waiting for deployment to be ready...${NC}"
kubectl rollout status deployment/dev-web-analyzer -n web-analyzer-dev --timeout=120s

# Get service info
echo ""
echo -e "${GREEN}✅ Deployment complete!${NC}"
echo ""
echo -e "${YELLOW}📍 To access the application:${NC}"
echo "   kubectl port-forward -n web-analyzer-dev svc/dev-web-analyzer 8080:80"
echo ""
echo "   Then open: http://localhost:8080"
echo ""

# Optionally start port-forward
read -p "Start port-forwarding now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}🌐 Starting port-forward... (Press Ctrl+C to stop)${NC}"
    kubectl port-forward -n web-analyzer-dev svc/dev-web-analyzer 8080:80
fi
