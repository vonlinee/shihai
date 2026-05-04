#!/bin/bash

# Start Backend Server (Linux)
# Usage: Run from project root directory

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_PATH="/opt/shihai"
BACKEND_PATH="$DEPLOY_PATH/shihai-server"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

if [ ! -f "$BACKEND_PATH" ]; then
    echo -e "${RED}Error: Backend executable not found at $BACKEND_PATH${NC}"
    echo -e "${YELLOW}Please run ./scripts/linux/deploy-linux.sh first${NC}"
    exit 1
fi

echo -e "${CYAN}Starting Shihai Backend Server...${NC}"
echo -e "${YELLOW}Server will run on http://localhost:8080${NC}"
echo -e "${GRAY}Press Ctrl+C to stop${NC}"
echo ""

cd "$DEPLOY_PATH"
./shihai-server
