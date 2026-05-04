#!/bin/bash

# Start Both Backend and Frontend (Linux)
# Usage: Run from project root directory

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$SCRIPTS_DIR")"
DEPLOY_PATH="/opt/shihai"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}Shihai Poetry Platform - Starting All${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# Check if deployment exists
if [ ! -d "$DEPLOY_PATH" ]; then
    echo -e "${RED}Error: Deployment not found at $DEPLOY_PATH${NC}"
    echo -e "${YELLOW}Please run ./scripts/linux/deploy-linux.sh first${NC}"
    exit 1
fi

# Function to cleanup on exit
cleanup() {
    echo ""
    echo -e "${YELLOW}Stopping services...${NC}"
    sudo pkill -f "shihai-server" 2>/dev/null || true
    sudo pkill -f "http.server" 2>/dev/null || true
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start Backend
echo -e "${GREEN}Starting Backend Server...${NC}"
cd "$DEPLOY_PATH"
sudo ./shihai-server &
BACKEND_PID=$!

# Wait for backend to start
sleep 2

# Start Frontend
echo -e "${GREEN}Starting Frontend Server...${NC}"
cd "$DEPLOY_PATH/frontend"
if command -v python3 &> /dev/null; then
    sudo python3 -m http.server 80 &
elif command -v python &> /dev/null; then
    sudo python -m http.server 80 &
fi
FRONTEND_PID=$!

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${GREEN}All services started!${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""
echo -e "Backend:  ${WHITE}http://localhost:8080${NC}"
echo -e "Frontend: ${WHITE}http://localhost${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# Wait for both processes
wait
