# 识海后端开发指南

本文档规定识海 Go 后端的项目结构、分层、接口、数据访问、认证和测试规则。通用 Go 代码要求见 [Go 开发规范](go-development-guidelines.md)，产品范围和系统现状分别见 [需求设计文档](requirements-design.md) 与 [系统设计文档](system-design.md)。

## 1. 项目概览

后端位于 `backend/`，使用 Go 1.21、Gin、GORM 和 PostgreSQL。当前服务入口为 `cmd/server/main.go`，API 前缀为 `/api`，默认监听 `8080` 端口。

### 1.1 核心目录

```text
backend/
├── cmd/server/          # 启动、依赖组装和路由注册
├── internal/config/     # 配置加载、数据库连接和 AutoMigrate
├── internal/database/   # 数据库初始化辅助能力
├── internal/dto/        # 请求和响应 DTO
├── internal/handlers/   # HTTP 适配层
├── internal/middleware/ # 认证、RBAC 等中间件
├── internal/models/     # GORM 模型和权限定义
├── internal/repository/ # 数据访问层
├── internal/services/   # 业务规则和流程编排
└── pkg/utils/           # 响应、JWT、Snowflake ID 等通用工具
```

### 1.2 本地运行

```bash
cd backend
go mod download
go run ./cmd/server
```

配置文件不是启动的强制条件。加载优先级和生产环境要求见“配置管理”章节。

### 1.3 常用验证

```bash
cd backend
gofmt -w <changed-go-files>
go test ./...
```

## 2. 分层边界

请求默认按以下方向调用：

```text
Route / Middleware -> Handler -> Service -> Repository -> GORM
```

### 2.1 Handler

Handler 负责：

- 从路径、查询、请求头和请求体绑定参数。
- 执行 HTTP 层格式校验。
- 从认证上下文读取当前用户。
- 调用 Service。
- 将业务结果映射为 HTTP 状态和统一响应。

Handler 不直接编写数据库查询，不承载可复用业务规则。健康检查等基础设施端点可以直接返回简单 JSON，但必须作为明确例外。

### 2.2 Service

Service 负责：

- 业务校验和状态转换。
- 多个 Repository 或其他 Service 的流程编排。
- 定义事务边界和幂等要求。
- 将持久化错误转换为调用方可判断的业务错误。

Service 不依赖 Gin Context，不直接构造 HTTP 响应。

### 2.3 Repository

Repository 负责：

- 封装 GORM 查询和写入。
- 处理分页、排序、预加载和数据库错误。
- 提供 Service 所需的持久化语义。

Repository 不决定用户是否有权限，也不拼装面向前端的展示文本。

### 2.4 DTO 与 Model

- 请求和响应使用 DTO，禁止直接把 GORM Model 作为稳定的外部契约。
- Model 描述持久化结构和关联关系。
- DTO 字段名、可空性、ID 表达和分页结构必须与前端契约一致。
- 一个 DTO 可以组合多个 Model 的结果，但组合逻辑应在 Service 或专用映射函数中完成。

## 3. 事务

单仓储、单条写入可以直接使用 Repository 方法。多个写操作必须原子完成时，由 Service 定义事务范围，并通过以下一种明确方式让相关 Repository 使用同一个事务：

- 事务协调器执行回调，并向回调传入绑定事务的 Repository 集合。
- Repository 提供 `WithTx(*gorm.DB)` 或等效方法，返回使用同一事务的新实例。
- Service 依赖定义好的 Unit of Work 接口。

禁止在一个流程中让部分 Repository 使用事务对象、其他 Repository 继续使用原始数据库连接。事务回调返回错误时必须回滚，提交失败也必须向上返回。

当前代码尚未形成统一事务抽象。新增跨仓储写流程时必须先选择并测试一种方案，不要复制彼此不兼容的事务写法。

## 4. API 契约

### 4.1 路径

- 当前业务 API 使用 `/api` 前缀，例如 `/api/poems`、`/api/admin/poems`。
- 资源路径使用小写和连字符，集合通常使用复数名词。
- 行为无法自然表达为 CRUD 时可以使用明确的动作子路径，例如 `/comments/vote`、`/text-conversion`。
- 当前没有 `/api/v1`。引入版本号必须作为跨端契约变更处理。

### 4.2 Swagger 文档

- 新增、修改或删除 HTTP 接口时，必须同步调整 `cmd/server/swagger_annotations_*.go` 中对应的 Swagger 注解。
- 注解必须与实际接口的请求方法、路径、认证要求、参数、请求体、响应结构和主要状态码保持一致。
- 修改注解后，必须在 `backend/` 目录重新生成 `docs/swagger/`，并运行 `go test ./cmd/server` 验证 Swagger UI 和文档接口。
- `docs/swagger/` 是由注解生成的文件，不得直接手工修改；使用以下命令重新生成：

```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g cmd/server/main.go -o docs/swagger --parseInternal
```

### 4.3 统一响应

普通业务接口使用：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

分页数据当前使用：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [],
    "total": 0,
    "page": 1,
    "pageSize": 10
  }
}
```

Handler 优先使用 `pkg/utils` 的响应函数。若需要新增响应形式，应先统一扩展响应工具和前端类型，避免各 Handler 自行定义包络。

### 4.4 分页

- 查询参数使用 `page` 和 `pageSize`，页码从 1 开始。
- DTO 应设置合理默认值和最大值；Repository 不能接受无限制的全表列表请求。
- 排序字段必须使用白名单映射，禁止把客户端字符串直接拼入 SQL。

### 4.5 Snowflake ID

- Model 中的业务 ID 使用 `uint64`。
- 对外响应 DTO 将 Snowflake ID 序列化为字符串。
- 请求 DTO 使用统一 `RequestID` 或等效解析类型，兼容字符串和必要的历史数字格式。
- JWT 中可被前端读取的用户 ID 同样使用字符串输出。
- 修改 ID 契约时必须补充序列化、反序列化和前端精度测试。

## 5. 错误与状态码

Service 应返回可分类的错误，Handler 根据语义映射：

| HTTP 状态 | 场景 |
| --- | --- |
| 400 | JSON、路径或查询参数格式错误 |
| 401 | 缺少、无效或过期的认证信息 |
| 403 | 身份有效但无权执行操作 |
| 404 | 指定资源不存在 |
| 409 | 唯一约束、重复操作或状态冲突 |
| 422 | 参数格式正确但不满足业务规则 |
| 500 | 未知系统错误、数据库或依赖故障 |

- 不通过错误文本字符串判断状态码；使用哨兵错误、自定义错误类型或明确错误码。
- 500 响应不得暴露 SQL、连接信息、文件路径或调用堆栈。
- 详细上下文写入服务端日志，对客户端返回稳定、可理解的消息。
- 现有部分 Handler 将业务错误统一返回 400，相关代码被修改时应逐步迁移到语义化映射。

## 6. 数据模型与数据库

### 6.1 模型

- 当前业务模型通常内嵌 `BaseModel`，获得 Snowflake ID、时间戳、软删除和审计用户字段。
- 表名通过 `TableName()` 明确指定为当前使用的单数形式。
- 关联字段应标明外键、索引和必要的唯一约束。
- 密码和内部敏感字段使用 `json:"-"`，关联对象按需使用 `omitempty`。
- GORM 标签中的 `comment` 用于数据库含义，Go 注释用于代码意图，两者不要求机械重复。

### 6.2 查询

- 使用 GORM 参数绑定，禁止拼接用户输入生成 SQL。
- 列表查询必须分页或有明确上限。
- 按需使用 `Preload`，避免 N+1 查询和无关的大对象加载。
- 新增高频查询条件时同步评估索引，并通过查询计划或基准数据验证。
- 批量删除、更新或查询必须处理空 ID 集合和最大批量限制。

### 6.3 软删除

包含 `gorm.DeletedAt` 的业务模型默认使用软删除。审计日志、关联表或有明确生命周期的数据可以采用不同策略，但必须在模型和业务规则中说明，不能只依赖团队记忆。

### 6.4 迁移

当前服务启动会执行 AutoMigrate，而 `database/migrations/001_init.sql` 保留了旧的复数表名和自增主键结构。修改模型前必须确认目标环境的建库方式，不能假设旧迁移与当前模型完全一致。

迁移策略统一前，不要在普通功能任务中同时修改模型、旧初始化脚本和生产数据；数据库结构变更应单独评审回滚和兼容方案。

## 7. 认证与权限

- JWT 从 `Authorization: Bearer <token>` 读取。
- 密码必须使用 bcrypt 等密码哈希算法，禁止明文保存。
- 路由注册处明确声明认证和权限中间件。
- 权限码统一定义在 `models/permission_codes.go`，不要在 Handler 中散落字符串字面量。
- 新增权限时同步更新权限元数据、默认角色映射、路由和权限测试。
- 前端权限控制不是安全边界；后端必须独立拒绝未授权请求。
- 高权限初始化不得依赖公开部署中的固定默认密码。

## 8. 配置管理

配置文件路径按命令行 `-config`、`CONFIG_FILE`、默认 `config.json` 的顺序选择。最终选中的文件不存在时保留内置开发默认值，不继续回退到其他文件；环境变量最后应用并覆盖已加载值。

当前支持的环境变量包括：

- `SERVER_PORT`、`GIN_MODE`。
- `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASSWORD`、`DB_NAME`、`DB_SSLMODE`。
- `REDIS_HOST`、`REDIS_PORT`、`REDIS_PASSWORD`、`REDIS_DB`，当前只有配置结构，未接入 Redis 客户端。
- `JWT_SECRET`、`JWT_EXPIRES_IN`。

生产环境必须覆盖开发数据库密码。当前配置加载器虽然读取 `JWT_SECRET` 和 `JWT_EXPIRES_IN`，但启动流程尚未将其传入 `pkg/utils` 的 JWT 实现；在完成接线和配置生效测试前，不得假设设置环境变量已经替换实际签名密钥和有效期。

密钥、令牌和生产连接信息不得提交到版本库或写入日志。

## 9. 日志

- 当前可以使用标准库 `log` 记录启动和运行信息，但标准库日志没有原生级别与结构化字段。
- 在引入统一日志组件前，不把 DEBUG、INFO、WARN、ERROR 写成已具备的强制分级能力。
- 错误日志至少包含请求方法、路径、业务模块和可关联的错误上下文。
- 禁止记录密码、完整 Token、数据库密码和不必要的个人信息。
- 禁止使用 `fmt.Println` 保留临时调试输出。
- 当前响应工具会打印错误堆栈，生产环境是否保留需要通过独立日志安全任务评估。

## 10. 测试

- Service 的重要业务规则应有单元测试。
- Repository 的查询、分页、约束和事务行为使用测试数据库或可验证的数据库替身测试。
- Handler 使用 `httptest` 覆盖参数错误、认证授权、状态码和响应结构。
- 权限、Snowflake ID、批量操作和文本转换等共享能力必须覆盖边界条件。
- 修复线上或可复现缺陷时，优先增加能够防止回归的测试。

```bash
cd backend
go test ./...
go test ./internal/services/... -v
```

## 11. Git

提交格式、检查流程和分支策略统一遵守 [Git 提交规范](git-commit-guidelines.md)，本文档不重复定义。
