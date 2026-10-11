# 跨平台脚本入口

使用 `scripts/shihai.mjs` 或 `scripts/shihai.py` 作为统一入口。两个入口都不依赖 Bash、PowerShell、Windows Terminal 或其他第三方终端程序，因此可以直接从 Windows 命令提示符、PowerShell、Linux/macOS 终端执行。

## 环境要求

- Node.js 18 或更高版本（建议使用项目当前的 Node.js LTS）。
- Python 3.9 或更高版本（仅使用 Python 标准库）。
- 部署或开发启动需要 Node.js、npm 和 Go 1.21 或更高版本。
- 后端运行仍需要 PostgreSQL，以及可用的 `backend/config.json` 或对应环境变量。

所有命令都从项目根目录执行，通用格式为：

```text
node scripts/shihai.mjs <命令> [选项]
```

Python 版本：

```text
python scripts/shihai.py <命令> [选项]
```

## 命令

### 一键启动

默认重新构建嵌入式部署产物，并在当前终端启动服务：

```text
node scripts/shihai.mjs start
```

开发模式会在当前终端管理前后端两个进程：

```text
node scripts/shihai.mjs start --mode development
node scripts/shihai.mjs start --mode development --backend-port 9090
```

独立前端部署：

```text
node scripts/shihai.mjs start --frontend-mode standalone --frontend-port 4173
```

使用已有部署产物：

```text
node scripts/shihai.mjs start --no-redeploy
```

预览启动计划，不构建也不启动服务：

```text
node scripts/shihai.mjs start --mode development --dry-run
```

启动计划和实际启动信息都会明确打印 `Frontend URL` 与 `Backend URL`。开发模式前端默认是 `http://localhost:3000`；独立部署模式分别使用 `--frontend-port` 和 `--backend-port`；嵌入式部署模式前后端都由后端端口提供服务。

后台运行、查看状态和停止：

```text
node scripts/shihai.mjs start --detach
node scripts/shihai.mjs start --status
node scripts/shihai.mjs start --stop
```

后台运行的 PID 和日志位于部署目录的 `run/start-all.pid`、`logs/start-all.log`；开发模式位于项目的 `tmp/universal-runtime/`。

### 单独部署

```text
node scripts/shihai.mjs deploy
node scripts/shihai.mjs deploy --frontend-mode embedded
node scripts/shihai.mjs deploy --frontend-mode embedded --build-only
node scripts/shihai.mjs deploy --deploy-path /srv/shihai
```

Windows 路径直接使用 Windows 格式，例如 `--deploy-path D:\shihai-deploy`。独立模式允许使用 `--skip-frontend` 或 `--skip-backend`；嵌入模式必须同时构建前后端。

### 单独启动服务

```text
node scripts/shihai.mjs backend --deploy-path /opt/shihai --port 9090
node scripts/shihai.mjs frontend --deploy-path /opt/shihai --port 4173
node scripts/shihai.mjs frontend --deploy-path /opt/shihai --port 4173 --dry-run
node scripts/shihai.mjs start-development-frontend --backend-port 9090
```

### 检查

```text
node scripts/shihai.mjs test
node scripts/shihai.mjs test --target backend
node scripts/shihai.mjs test --target frontend
```

完整选项可执行 `node scripts/shihai.mjs --help` 查看。原有 `scripts/linux` 和 `scripts/windows` 脚本保持兼容；新项目或跨平台自动化场景建议统一使用本入口。

Python 入口的完整选项可执行 `python scripts/shihai.py --help` 查看。两种入口的命令和参数保持一致，可以按环境已有的运行时选择。

## Python 版本测试

```text
python -m unittest discover -s scripts/tests -p "test_shihai.py"
```
