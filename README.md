



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

# Optional development proxy target when the backend uses a custom port
VITE_BACKEND_URL=http://localhost:8080

# Production
VITE_API_BASE_URL=https://your-domain.com/api
```
## Backend:

### Development:
```bash
cd backend
go mod tidy
go run ./cmd/server

# Override the backend port (takes precedence over SERVER_PORT and config.json)
go run ./cmd/server -port 9090
```

### Build & Deploy:

#### Build for current platform:
```bash
cd backend
go mod tidy
go build -o shihai-server.exe ./cmd/server
```

#### Build for Linux (AMD64):
```bash
cd backend
GOOS=linux GOARCH=amd64 go build -o shihai-server-linux ./cmd/server
```

#### Build for Windows (AMD64):
```bash
cd backend
GOOS=windows GOARCH=amd64 go build -o shihai-server.exe ./cmd/server
```

#### Build for macOS (AMD64):
```bash
cd backend
GOOS=darwin GOARCH=amd64 go build -o shihai-server-mac ./cmd/server
```

#### Build for macOS (ARM64/M1):
```bash
cd backend
GOOS=darwin GOARCH=arm64 go build -o shihai-server-mac-arm ./cmd/server
```

#### Run the built binary:
```bash
# Windows
./shihai-server.exe

# Custom port
./shihai-server.exe -port 9090

# Linux/Mac
./shihai-server-linux
# or
./shihai-server-mac
```

## Database:
Execute database/migrations/001_init.sql in PostgreSQL

## Swagger API Documentation:

Swagger metadata and endpoint annotations are maintained in `backend/internal/apidocs/`. Generate the Swagger package from the `backend/` directory:

```bash
cd backend
go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g internal/apidocs/doc.go -o docs/swagger --parseInternal
```

The generated files are written to `backend/docs/swagger/` and should not be edited manually.

## Cross-platform Scripts:

The repository provides a platform-independent launcher that can be run from
Windows Command Prompt, PowerShell, Linux, or macOS terminals. It uses only
the selected runtime's standard library and invokes the project tools directly.

Node.js version:

```text
node scripts/shihai.mjs start
node scripts/shihai.mjs start --mode development
node scripts/shihai.mjs start --frontend-mode standalone --frontend-port 4173
```

Python version:

```text
python scripts/shihai.py start
python scripts/shihai.py start --mode development
python scripts/shihai.py start --frontend-mode standalone --frontend-port 4173
```

Common operations are available in both versions:

```text
# Preview the startup plan and access URLs
node scripts/shihai.mjs start --mode development --backend-port 9090 --dry-run

# Build and deploy without starting services
node scripts/shihai.mjs deploy --frontend-mode embedded

# Run in the background, inspect status, and stop
node scripts/shihai.mjs start --detach
node scripts/shihai.mjs start --status
node scripts/shihai.mjs start --stop

# Run project checks
node scripts/shihai.mjs test
python scripts/shihai.py test
```

The launcher prints `Frontend URL` and `Backend URL` for every startup plan.
The Python commands are direct replacements for the Node.js commands, for
example `python scripts/shihai.py start --status` and
`python scripts/shihai.py deploy --frontend-mode embedded`.

See [cross-platform scripts documentation](scripts/README.md) for all options,
deployment modes, standalone service commands, and runtime requirements.

## Quick Deploy Scripts:

### Windows:
```powershell
# Default: rebuild and start the backend with the embedded frontend
.\scripts\windows\start-all.ps1

# Development mode: open frontend and backend in two Windows Terminal tabs
.\scripts\windows\start-all.ps1 -Mode Development

# Use a custom backend port; the development proxy follows automatically
.\scripts\windows\start-all.ps1 -Mode Development -BackendPort 9090

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

# Use a custom backend port; the development proxy follows automatically
./scripts/linux/start-all.sh --mode development --backend-port 9090

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


