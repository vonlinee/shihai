#!/bin/bash

# Shihai Poetry Platform - Linux Deployment Script
# This script builds and deploys both frontend and backend locally
# Usage: Run from project root directory

set -e

# Get script directory and project root (script is in scripts/linux/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$SCRIPTS_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
BACKEND_PORT="8080"
FRONTEND_DIST_PATH="$PROJECT_ROOT/frontend/dist"
BACKEND_BUILD_PATH="$PROJECT_ROOT/backend/shihai-server"
DEPLOY_PATH="/opt/shihai"

# Flags
BUILD_ONLY=false
SKIP_FRONTEND=false
SKIP_BACKEND=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --build-only)
            BUILD_ONLY=true
            shift
            ;;
        --skip-frontend)
            SKIP_FRONTEND=true
            shift
            ;;
        --skip-backend)
            SKIP_BACKEND=true
            shift
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}Shihai Poetry Platform Deployment Tool${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

if [ "$SKIP_FRONTEND" = false ]; then
    if ! command_exists node; then
        echo -e "${RED}Error: Node.js is not installed or not in PATH${NC}"
        exit 1
    fi
    echo -e "  ${GREEN}[OK]${NC} Node.js found"
fi

if [ "$SKIP_BACKEND" = false ]; then
    if ! command_exists go; then
        echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
        exit 1
    fi
    echo -e "  ${GREEN}[OK]${NC} Go found"
fi

echo ""

# Build Frontend
if [ "$SKIP_FRONTEND" = false ]; then
    echo -e "${YELLOW}Building Frontend...${NC}"
    cd "$PROJECT_ROOT/frontend"
    
    echo -e "  ${GRAY}Installing dependencies...${NC}"
    npm install
    
    echo -e "  ${GRAY}Building production bundle...${NC}"
    npm run build
    
    cd ..
    echo -e "  ${GREEN}[OK]${NC} Frontend built successfully"
    echo ""
fi

# Build Backend
if [ "$SKIP_BACKEND" = false ]; then
    echo -e "${YELLOW}Building Backend...${NC}"
    cd "$PROJECT_ROOT/backend"
    
    echo -e "  ${GRAY}Downloading dependencies...${NC}"
    go mod tidy
    
    echo -e "  ${GRAY}Building executable...${NC}"
    go build -o shihai-server cmd/server/main.go
    
    cd ..
    echo -e "  ${GREEN}[OK]${NC} Backend built successfully"
    echo ""
fi

# Deploy
if [ "$BUILD_ONLY" = false ]; then
    echo -e "${YELLOW}Deploying to local machine...${NC}"
    
    # Create deployment directory (may require sudo)
    if [ ! -d "$DEPLOY_PATH" ]; then
        echo -e "  ${GRAY}Creating deployment directory...${NC}"
        sudo mkdir -p "$DEPLOY_PATH/frontend"
        sudo chown -R $(whoami):$(whoami) "$DEPLOY_PATH"
    fi
    
    # Copy frontend files
    if [ "$SKIP_FRONTEND" = false ]; then
        echo -e "  ${GRAY}Copying frontend files...${NC}"
        cp -r "$FRONTEND_DIST_PATH/"* "$DEPLOY_PATH/frontend/"
        echo -e "  ${GREEN}[OK]${NC} Frontend deployed"
    fi
    
    # Copy backend files
    if [ "$SKIP_BACKEND" = false ]; then
        echo -e "  ${GRAY}Copying backend executable...${NC}"
        cp "$BACKEND_BUILD_PATH" "$DEPLOY_PATH/"
        chmod +x "$DEPLOY_PATH/shihai-server"
        echo -e "  ${GREEN}[OK]${NC} Backend deployed"
    fi
    
    echo ""
    echo -e "${CYAN}========================================${NC}"
    echo -e "${GREEN}Deployment Complete!${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo ""
    echo -e "Deployment location: ${WHITE}$DEPLOY_PATH${NC}"
    echo ""
    echo -e "${YELLOW}To start the application:${NC}"
    echo -e "  1. Start Backend:  sudo ./scripts/linux/start-backend.sh"
    echo -e "  2. Start Frontend: sudo ./scripts/linux/start-frontend.sh"
    echo -e "  3. Or start both:  sudo ./scripts/linux/start-all.sh"
    echo ""
else
    echo -e "${YELLOW}Build only mode - skipping deployment${NC}"
fi
