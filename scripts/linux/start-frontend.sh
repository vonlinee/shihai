#!/bin/bash

# Start Frontend Server (Linux)
# Usage: Run from project root directory

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_PATH="/opt/shihai/frontend"
PORT="80"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

if [ ! -d "$DEPLOY_PATH" ]; then
    echo -e "${RED}Error: Frontend files not found at $DEPLOY_PATH${NC}"
    echo -e "${YELLOW}Please run ./scripts/linux/deploy-linux.sh first${NC}"
    exit 1
fi

echo -e "${CYAN}Starting Shihai Frontend Server...${NC}"
echo -e "${YELLOW}Server will run on http://localhost:$PORT${NC}"
echo -e "${GRAY}Press Ctrl+C to stop${NC}"
echo ""

cd "$DEPLOY_PATH"

# Try Python first, then fallback to other methods
if command -v python3 &> /dev/null; then
    sudo python3 -m http.server $PORT
elif command -v python &> /dev/null; then
    sudo python -m http.server $PORT
elif command -v nginx &> /dev/null; then
    echo -e "${YELLOW}Using nginx to serve frontend...${NC}"
    sudo nginx -c "$DEPLOY_PATH/nginx.conf"
else
    echo -e "${RED}Error: No suitable web server found${NC}"
    echo -e "${YELLOW}Please install Python or nginx${NC}"
    exit 1
fi
