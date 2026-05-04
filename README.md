



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
# Build and deploy everything
.\scripts\windows\deploy-windows.ps1

# Or build only (skip deployment)
.\scripts\windows\deploy-windows.ps1 -BuildOnly

# Skip frontend or backend
.\scripts\windows\deploy-windows.ps1 -SkipFrontend
.\scripts\windows\deploy-windows.ps1 -SkipBackend

# Start services
.\scripts\windows\start-all.ps1        # Start both
.\scripts\windows\start-backend.ps1    # Start backend only
.\scripts\windows\start-frontend.ps1   # Start frontend only
```

### Linux:
```bash
# Make scripts executable first
chmod +x scripts/linux/*.sh

# Build and deploy everything
sudo ./scripts/linux/deploy-linux.sh

# Or build only (skip deployment)
sudo ./scripts/linux/deploy-linux.sh --build-only

# Skip frontend or backend
sudo ./scripts/linux/deploy-linux.sh --skip-frontend
sudo ./scripts/linux/deploy-linux.sh --skip-backend

# Start services
sudo ./scripts/linux/start-all.sh        # Start both
sudo ./scripts/linux/start-backend.sh    # Start backend only
sudo ./scripts/linux/start-frontend.sh   # Start frontend only
```


