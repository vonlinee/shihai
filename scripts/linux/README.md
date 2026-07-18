# Linux 启动与部署脚本

本目录提供识海项目在 Linux 环境下的开发启动、本地部署和服务启动脚本。

## 环境要求

- Bash 4.3 或更高版本。
- Node.js 和 npm，用于前端开发与构建。
- Go 1.21 或更高版本，用于后端开发与构建。
- PostgreSQL，以及可用的 `backend/config.json` 或对应环境变量。
- 默认部署目录为 `/opt/shihai`，首次创建或权限不足时需要可用的 `sudo`。

首次使用前添加执行权限：

```bash
chmod +x scripts/linux/*.sh scripts/linux/tests/*.sh
```

## 一键启动

统一入口为 `start-all.sh`。默认重新构建嵌入式部署产物，并在当前终端启动一个包含前端资源的 Go 服务：

```bash
./scripts/linux/start-all.sh
```

Linux 终端没有统一的 Tab 控制协议，因此脚本不会依赖 GNOME Terminal、Konsole 或 tmux。开发模式和独立部署模式会在当前 shell 管理两个子进程；按 `Ctrl+C` 或任一进程退出时，脚本会停止其余进程。

### 开发模式

```bash
./scripts/linux/start-all.sh --mode development
```

| 进程 | 命令 | 地址 |
| --- | --- | --- |
| frontend | `npm run dev` | `http://localhost:3000` |
| backend | `go run ./cmd/server` | `http://localhost:8080` |

前端开发服务器会把 `/api` 请求代理到后端的 `8080` 端口。

### 独立前端部署模式

```bash
./scripts/linux/start-all.sh --mode deployment --frontend-mode standalone
```

该模式默认重新构建并部署到 `/opt/shihai`，然后在当前终端管理前后端两个进程：

- 前端使用项目已有的 Vite Preview 提供静态文件和 SPA 路由回退，默认端口为 `80`。
- 后端启动 `/opt/shihai/shihai-server`，默认地址为 `http://localhost:8080`。
- 前端构建时默认把 API 地址设置为 `http://localhost:8080/api`。

Linux 普通用户通常不能监听 `80` 端口，可以改用高端口：

```bash
./scripts/linux/start-all.sh \
  --mode deployment \
  --frontend-mode standalone \
  --frontend-port 4173
```

部署目录结构：

```text
/opt/shihai/
├── frontend/
├── frontend-mode.txt
├── config.json
└── shihai-server
```

### 前端嵌入部署模式

无参数启动即为嵌入部署模式，完整写法为：

```bash
./scripts/linux/start-all.sh --mode deployment --frontend-mode embedded
```

脚本先构建前端，再使用 `embed_frontend` 构建标签把资源嵌入 `shihai-server`。前端页面和 API 都由 `http://localhost:8080` 提供。临时嵌入目录会在构建完成或失败后自动清理。

### 使用已有部署产物

部署模式默认重新构建。跳过构建并直接使用现有产物：

```bash
# 复用嵌入式部署产物
./scripts/linux/start-all.sh --no-redeploy

# 复用独立部署产物
./scripts/linux/start-all.sh \
  --frontend-mode standalone \
  --no-redeploy \
  --frontend-port 4173
```

脚本会读取 `frontend-mode.txt`，检查已有产物是否与请求模式一致。

### 后台启动与停止

`start-all.sh` 可以独立完成后台启动、状态查询和停止，不需要手工组合 `nohup`、日志重定向和 PID 文件：

```bash
# 默认嵌入模式：前台完成构建部署，然后转入后台
./scripts/linux/start-all.sh --detach

# 查看状态
./scripts/linux/start-all.sh --status

# 查看日志
tail -f /opt/shihai/logs/start-all.log

# 停止脚本管理的全部服务
./scripts/linux/start-all.sh --stop
```

部署模式的 PID 文件为 `/opt/shihai/run/start-all.pid`。使用 `--deploy-path` 时，日志和 PID 会写入对应部署目录；查询状态和停止时需要传入相同的 `--deploy-path`。

开发模式后台运行：

```bash
./scripts/linux/start-all.sh --mode development --detach
./scripts/linux/start-all.sh --mode development --status
./scripts/linux/start-all.sh --mode development --stop
```

开发模式的运行文件位于项目 `tmp/linux-runtime/`。后台启动前会检查是否已有受管理实例，避免重复启动。

### 自定义部署配置

```bash
# 自定义部署目录
./scripts/linux/start-all.sh --deploy-path /srv/shihai

# 自定义独立前端的 API 地址
./scripts/linux/start-all.sh \
  --frontend-mode standalone \
  --standalone-api-base-url https://api.example.com/api
```

### 预览启动计划

`--dry-run` 只显示模式、进程、工作目录和命令，不构建、不部署，也不启动服务：

```bash
./scripts/linux/start-all.sh --dry-run
./scripts/linux/start-all.sh --mode development --dry-run
```

## 启动参数

| 参数 | 可选值 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--mode` | `development`、`deployment` | `deployment` | 选择开发或部署启动模式 |
| `--frontend-mode` | `standalone`、`embedded` | `embedded` | 选择部署模式下的前端产物形式 |
| `--redeploy` | 开关 | 开启 | 强制部署启动前重新构建 |
| `--no-redeploy` | 开关 | 关闭 | 使用已有部署产物 |
| `--deploy-path` | Linux 绝对路径 | `/opt/shihai` | 部署产物目录 |
| `--standalone-api-base-url` | URL | `http://localhost:8080/api` | 独立前端构建使用的 API 地址 |
| `--frontend-port` | `1` 至 `65535` | `80` | 独立前端运行端口 |
| `--detach` | 开关 | 关闭 | 完成必要部署后在后台启动服务 |
| `--status` | 开关 | 关闭 | 查看后台实例状态 |
| `--stop` | 开关 | 关闭 | 停止后台实例及其子进程 |
| `--dry-run` | 开关 | 关闭 | 仅预览启动计划 |

## 单独部署

```bash
# 独立前端模式
./scripts/linux/deploy-linux.sh

# 嵌入前端模式
./scripts/linux/deploy-linux.sh --frontend-mode embedded

# 只构建，不复制到部署目录
./scripts/linux/deploy-linux.sh --frontend-mode embedded --build-only

# 自定义部署目录
./scripts/linux/deploy-linux.sh --deploy-path /srv/shihai
```

`--skip-frontend` 和 `--skip-backend` 仅适用于独立前端模式。嵌入模式需要同时构建前后端。

## 单独启动部署服务

```bash
./scripts/linux/start-backend.sh --deploy-path /opt/shihai
./scripts/linux/start-frontend.sh --deploy-path /opt/shihai --port 4173
```

前端脚本还支持 `--dry-run`。在当前终端按 `Ctrl+C` 即可停止服务。
