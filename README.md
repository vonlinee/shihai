



# To Run the Project

## Frontend:

### Development:
```shell
cd frontend
npm install
npm run dev
```

### Build & Deploy:

#### Build for production:
```bash
cd frontend
npm install
npm run build
```

The build output will be in the `frontend/dist` directory.

#### Preview production build locally:
```bash
cd frontend
npm run preview
```

#### Deploy to production server:

After building, copy the `dist` folder contents to your web server:

```bash
# Example: Copy to nginx web root
sudo cp -r frontend/dist/* /var/www/html/

# Or use rsync for remote deployment
rsync -avz frontend/dist/ user@server:/var/www/shihai/
```

#### Environment Variables:

Create a `.env` file in the frontend directory for different environments:

```env
# Development
VITE_API_BASE_URL=http://localhost:8080/api

# Production
VITE_API_BASE_URL=https://your-domain.com/api
```
## Backend:

### Development:
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### Build & Deploy:

#### Build for current platform:
```bash
cd backend
go mod tidy
go build -o shihai-server.exe cmd/server/main.go
```

#### Build for Linux (AMD64):
```bash
cd backend
GOOS=linux GOARCH=amd64 go build -o shihai-server-linux cmd/server/main.go
```

#### Build for Windows (AMD64):
```bash
cd backend
GOOS=windows GOARCH=amd64 go build -o shihai-server.exe cmd/server/main.go
```

#### Build for macOS (AMD64):
```bash
cd backend
GOOS=darwin GOARCH=amd64 go build -o shihai-server-mac cmd/server/main.go
```

#### Build for macOS (ARM64/M1):
```bash
cd backend
GOOS=darwin GOARCH=arm64 go build -o shihai-server-mac-arm cmd/server/main.go
```

#### Run the built binary:
```bash
# Windows
./shihai-server.exe

# Linux/Mac
./shihai-server-linux
# or
./shihai-server-mac
```

## Database:
Execute database/migrations/001_init.sql in PostgreSQL

## Quick Deploy Scripts:

### Windows:
```powershell
# Default: rebuild and start the backend with the embedded frontend
.\scripts\windows\start-all.ps1

# Development mode: open frontend and backend in two Windows Terminal tabs
.\scripts\windows\start-all.ps1 -Mode Development

# Standalone deployment mode: rebuild by default, then open two tabs
.\scripts\windows\start-all.ps1 -Mode Deployment -FrontendMode Standalone

# Reuse existing deployment artifacts
.\scripts\windows\start-all.ps1 -Redeploy:$false

# Build and deploy everything
.\scripts\windows\deploy-windows.ps1

# Build and deploy with the frontend embedded in the Go executable
.\scripts\windows\deploy-windows.ps1 -FrontendMode Embedded

# Or build only (skip deployment)
.\scripts\windows\deploy-windows.ps1 -BuildOnly

# Skip frontend or backend
.\scripts\windows\deploy-windows.ps1 -SkipFrontend
.\scripts\windows\deploy-windows.ps1 -SkipBackend

# Start deployed services individually
.\scripts\windows\start-backend.ps1    # Start backend only
.\scripts\windows\start-frontend.ps1   # Start frontend only
```

See [Windows scripts documentation](scripts/windows/README.md) for all modes and parameters.

### Linux:
```bash
# Make scripts executable first
chmod +x scripts/linux/*.sh scripts/linux/tests/*.sh

# Default: rebuild and start the backend with the embedded frontend
./scripts/linux/start-all.sh

# Development mode: manage frontend and backend in the current terminal
./scripts/linux/start-all.sh --mode development

# Standalone deployment mode
./scripts/linux/start-all.sh --frontend-mode standalone --frontend-port 4173

# Reuse existing deployment artifacts
./scripts/linux/start-all.sh --no-redeploy

# Start in background, inspect status, and stop
./scripts/linux/start-all.sh --detach
./scripts/linux/start-all.sh --status
./scripts/linux/start-all.sh --stop

# Build and deploy everything
./scripts/linux/deploy-linux.sh

# Build and deploy with the frontend embedded in the Go executable
./scripts/linux/deploy-linux.sh --frontend-mode embedded

# Or build only (skip deployment)
./scripts/linux/deploy-linux.sh --build-only

# Skip frontend or backend
./scripts/linux/deploy-linux.sh --skip-frontend
./scripts/linux/deploy-linux.sh --skip-backend

# Start deployed services individually
./scripts/linux/start-backend.sh
./scripts/linux/start-frontend.sh --port 4173
```

See [Linux scripts documentation](scripts/linux/README.md) for all modes and parameters.


