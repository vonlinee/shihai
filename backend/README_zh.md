# 识海 (Shihai) 后端项目

这是“识海”（一个学习古诗词的平台）的后端服务实现。本项目采用 Go（Golang）语言开发，使用 Gin Web 框架处理路由，并使用 GORM 与 PostgreSQL 进行数据库交互。

## 目录结构

后端项目遵循 Go 标准项目布局（Standard Go Project Layout）规范。该架构将业务逻辑、HTTP路由以及数据访问清晰地分层解耦。

```text
backend/
├── cmd/                        # 项目的主应用程序入口
│   └── server/                 # 包含启动服务器的主要入口点 (main.go)
│
├── internal/                   # 私有应用程序代码与核心库（不暴露给外部使用）
│   ├── config/                 # 系统环境配置加载和数据库初始化
│   ├── dto/                    # 数据传输对象 (DTO，用于API输入输出数据的结构体和校验逻辑)
│   ├── handlers/               # HTTP 请求处理器 (Controller 层，负责接收请求与返回响应)
│   ├── middleware/             # Gin 中间件 (包含认证、RBAC权限控制、CORS跨域等)
│   ├── models/                 # 数据库 Schema 定义和领域模型
│   ├── repository/             # 数据访问层 (DAL，封装具体的数据库读写与GORM逻辑)
│   └── services/               # 核心业务逻辑层 (协调数据交互，处理核心规则)
│
├── pkg/                        # 可供外部或本项目其它模块复用的公共类库
│   └── utils/                  # 各种通用工具函数 (如 JWT 解析、密码哈希、雪花算法等)
│
├── .idea/                      # JetBrains IDE 配置文件 (仅限本地开发环境)
├── bin/                        # 编译后生成的可执行文件目录
├── Golang开发规范.md           # Go 语言项目开发规范及指导手册
├── config.json                 # 当前活动的配置文件
├── config.example.json         # 配置文件示例
├── go.mod                      # Go module 依赖文件
└── go.sum                      # Go module 校验和文件
```

## 架构层级概述

- **表示层 / 控制器 (`internal/handlers`)**: 负责对外提供 RESTful API。利用 Gin 框架解析 HTTP 请求、提取参数并使用 `dto` 验证数据格式，随后将请求转交至 `services` 层，最后返回标准 JSON 响应包。
- **业务逻辑层 (`internal/services`)**: API应用的大脑，承载和编排系统的核心运转逻辑和事务流程。它在 `handlers` 与 `repositories` 之间起着承上启下的作用，保证业务与特定的 HTTP 格式或者特定的数据库相互隔离。
- **数据访问层 (`internal/repository`)**: 处理所有的数据库通信。该层专注执行经过封装的数据库查询/GORM命令，并将底层数据准确地读取与转换至对应的数据结构（`models`）中供服务层调用。

## 环境配置与运行

1. **安装依赖项目**
   ```bash
   go mod download
   ```

2. **配置数据库**
   请复制 `config.example.json` 并重命名为 `config.json`。修改里面的配置以确保与您的 PostgreSQL 实例所需凭据相匹配。

3. **启动应用程序**
   
   ```bash
   go run cmd/server/main.go
   ```
   *注意：项目启动时，GORM的 auto migration 工具会自动执行以同步最新的数据表结构至 `internal/models/` 里面的定义状态*。
