# Windows 启动与部署脚本

本目录提供识海项目在 Windows 环境下的开发启动、本地部署和服务启动脚本。

## 环境要求

- Windows Terminal，命令 `wt.exe` 必须可用。
- Node.js 和 npm，用于前端开发与构建。
- Go 1.21 或更高版本，用于后端开发与构建。
- PostgreSQL，以及可用的 `backend/config.json` 或对应环境变量。

如果 PowerShell 禁止执行本地脚本，可以在当前终端执行：

```powershell
Set-ExecutionPolicy -Scope Process Bypass
```

## 一键启动

统一入口为 `start-all.ps1`。默认重新构建嵌入式部署产物，并启动一个包含前端资源的后端 Tab。

- 在 Windows Terminal 中执行脚本时，新 Tab 会通过 `wt.exe -w 0` 添加到当前活动窗口。
- 在普通 PowerShell、命令提示符或其他终端中执行时，脚本会创建一个新的 Windows Terminal 窗口。

### 开发模式

```powershell
.\scripts\windows\start-all.ps1 -Mode Development
```

| Tab | 命令 | 地址 |
| --- | --- | --- |
| Shihai Frontend | `npm run dev` | `http://localhost:3000` |
| Shihai Backend | `go run ./cmd/server` | `http://localhost:8080` |

前端开发服务器会把 `/api` 请求代理到后端的 `8080` 端口。

使用其他后端端口时，脚本会同时调整后端命令和 Vite 开发代理：

```powershell
.\scripts\windows\start-all.ps1 -Mode Development -BackendPort 9090
```

### 独立前端部署模式

```powershell
.\scripts\windows\start-all.ps1 -Mode Deployment -FrontendMode Standalone
```

该模式默认重新构建并部署到 `C:\shihai-deploy`，然后打开两个 Tab：

- 前端 Tab 使用项目已有的 Vite Preview 提供 `frontend` 静态文件和 SPA 路由回退，默认地址为 `http://localhost`。
- 后端 Tab 启动 `shihai-server.exe`，默认地址为 `http://localhost:8080`。
- 前端构建时默认根据 `-BackendPort` 设置 API 地址，例如端口 `9090` 对应 `http://localhost:9090/api`；可通过 `-StandaloneApiBaseUrl` 显式覆盖。

部署目录结构：

```text
C:\shihai-deploy\
├── frontend\
├── frontend-mode.txt
└── shihai-server.exe
```

### 前端嵌入部署模式

```powershell
.\scripts\windows\start-all.ps1
```

完整写法为 `.\scripts\windows\start-all.ps1 -Mode Deployment -FrontendMode Embedded`。

该模式先构建前端，再使用 `embed_frontend` 构建标签把前端资源嵌入 `shihai-server.exe`。启动时只打开一个 `Backend (with embedded frontend)` Tab，前端页面和 API 都由 `http://localhost:8080` 提供。

嵌入构建使用的临时前端目录会在构建完成或失败后自动清理。

### 使用已有部署产物

部署模式默认重新构建部署。要跳过构建并直接使用现有产物：

```powershell
.\scripts\windows\start-all.ps1 -Redeploy:$false
.\scripts\windows\start-all.ps1 -Mode Deployment -FrontendMode Standalone -Redeploy:$false
```

脚本会读取部署目录中的 `frontend-mode.txt`，检查已有产物是否与请求的前端模式一致。

### 自定义部署目录

```powershell
.\scripts\windows\start-all.ps1 -Mode Deployment -DeployPath "D:\shihai-deploy"
```

所有模式都可以通过 `-BackendPort` 指定后端监听端口：

```powershell
.\scripts\windows\start-all.ps1 -BackendPort 9090
```

独立部署的后端 API 不在本机默认地址时，可以指定构建期 API 地址：

```powershell
.\scripts\windows\start-all.ps1 -Mode Deployment -FrontendMode Standalone -StandaloneApiBaseUrl "https://api.example.com/api"
```

### 自定义 Tab 名称

开发模式和独立部署模式可以分别指定前端、后端 Tab 标题：

```powershell
.\scripts\windows\start-all.ps1 -Mode Development -FrontendTabTitle "诗海前端" -BackendTabTitle "诗海后端"
```

嵌入部署模式只有一个 Tab，使用 `-BackendTabTitle` 设置标题：

```powershell
.\scripts\windows\start-all.ps1 -Mode Deployment -FrontendMode Embedded -BackendTabTitle "诗海应用"
```

### 预览启动计划

`-DryRun` 只显示模式、Tab、工作目录和命令，不构建、不部署，也不打开 Windows Terminal：

```powershell
.\scripts\windows\start-all.ps1 -DryRun
.\scripts\windows\start-all.ps1 -Mode Development -DryRun
```

## 启动参数

| 参数 | 可选值 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `-Mode` | `Development`、`Deployment` | `Deployment` | 选择开发或部署启动模式 |
| `-FrontendMode` | `Standalone`、`Embedded` | `Embedded` | 选择部署模式下的前端产物形式 |
| `-Redeploy` | `$true`、`$false` | `$true` | 部署启动前是否重新构建部署 |
| `-DeployPath` | Windows 路径 | `C:\shihai-deploy` | 部署产物目录 |
| `-BackendPort` | `1` 至 `65535` | `8080` | 后端监听端口；开发代理和独立前端默认 API 地址随之调整 |
| `-StandaloneApiBaseUrl` | URL | 根据 `-BackendPort` 生成 | 独立前端构建使用的 API 地址；显式设置时覆盖自动生成值 |
| `-FrontendTabTitle` | 文本 | `Shihai Frontend` | 开发或独立部署模式的前端 Tab 标题 |
| `-BackendTabTitle` | 文本 | 按模式确定 | 后端 Tab 标题；嵌入模式默认标注包含前端 |
| `-DryRun` | 开关 | 关闭 | 仅预览启动计划 |

## 单独部署

不启动服务，只构建并部署：

```powershell
# 独立前端模式
.\scripts\windows\deploy-windows.ps1

# 嵌入前端模式
.\scripts\windows\deploy-windows.ps1 -FrontendMode Embedded

# 只构建，不复制到部署目录
.\scripts\windows\deploy-windows.ps1 -BuildOnly

# 自定义部署目录
.\scripts\windows\deploy-windows.ps1 -DeployPath "D:\shihai-deploy"
```

`-SkipFrontend` 和 `-SkipBackend` 仅适用于独立前端模式。嵌入模式需要同时构建前后端，不能使用这两个参数。

## 单独启动部署服务

```powershell
.\scripts\windows\start-backend.ps1 -Port 9090
.\scripts\windows\start-frontend.ps1
```

两个脚本都支持 `-DeployPath`，前端脚本还支持 `-Port` 和 `-DryRun`。

在对应 Tab 中按 `Ctrl+C`，或关闭 Tab，即可停止服务。
