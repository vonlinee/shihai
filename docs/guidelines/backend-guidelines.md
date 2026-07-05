# 识海（shihai）Go Web 项目开发规范文档

> 本文档规定了识海古诗词学习平台后端 Go 项目的开发规范，适用于所有参与后端开发的成员，旨在统一代码风格、降低维护成本、提升代码质量。

---

## 项目概览

识海后端基于 Go、Gin、GORM 和 PostgreSQL 构建，采用分层架构组织 HTTP 处理、业务逻辑和数据访问。

### 技术栈

- Go
- Gin
- GORM
- PostgreSQL

### 核心目录

```text
backend/
├── cmd/server/          # 服务启动入口
├── internal/config/     # 配置加载和数据库初始化
├── internal/dto/        # API 请求和响应 DTO
├── internal/handlers/   # HTTP 处理器
├── internal/middleware/ # 鉴权、权限、CORS 等中间件
├── internal/models/     # GORM 模型和领域模型
├── internal/repository/ # 数据访问层
├── internal/services/   # 业务逻辑层
└── pkg/utils/           # 通用工具
```

### 分层约定

- Handler 层只处理 HTTP 输入输出、参数绑定、状态码和响应格式。
- Service 层承载业务规则、流程编排和事务边界。
- Repository 层封装数据库读写，不承载业务判断。
- DTO 负责接口输入输出结构，不直接复用数据库模型作为外部契约。
- Model 负责数据库结构、关联关系和持久化字段定义。

### 本地运行

```bash
cd backend
go mod download
go run cmd/server/main.go
```

启动前复制 `config.example.json` 为 `config.json`，并配置 PostgreSQL 连接信息。

### 常用验证

```bash
cd backend
gofmt -w <changed-go-files>
go test ./...
```

涉及数据库模型、服务层、仓储层或权限逻辑时，优先补充对应单元测试。

### 修改前阅读

- 新增或调整后端功能：阅读本文档和 [Go 开发规范](go-development-guidelines.md)。
- 调整系统模块、权限或数据关系：额外阅读 [系统设计文档](system-design.md)。
- 调整业务流程：额外阅读 [需求设计文档](requirements-design.md)。

---

## 目录

1. [项目结构规范](#一项目结构规范)
2. [命名规范](#二命名规范)
3. [代码风格规范](#三代码风格规范)
4. [分层架构规范](#四分层架构规范)
5. [API 设计规范](#五api-设计规范)
6. [错误处理规范](#六错误处理规范)
7. [数据模型规范](#七数据模型规范)
8. [数据库操作规范](#八数据库操作规范)
9. [认证与权限规范](#九认证与权限规范)
10. [日志规范](#十日志规范)
11. [配置管理规范](#十一配置管理规范)
12. [Git 提交与分支管理](#十二git-提交与分支管理)
13. [测试规范](#十三测试规范)

---

## 一、项目结构规范

项目采用标准 Go 项目布局，目录结构如下：

```
backend/
├── cmd/
│   └── server/          # 程序入口，只负责启动和初始化，不含业务逻辑
│       └── main.go
├── internal/            # 项目内部私有代码，不对外暴露
│   ├── config/          # 配置结构定义与加载
│   ├── dto/             # 请求/响应数据传输对象（Data Transfer Object）
│   ├── handlers/        # HTTP 请求处理层（Controller 层）
│   ├── middleware/      # Gin 中间件
│   ├── models/          # 数据库模型定义（与表结构对应）
│   ├── repository/      # 数据访问层（Repository 层）
│   └── services/        # 业务逻辑层（Service 层）
├── pkg/                 # 可复用的公共工具包
│   └── utils/           # 通用工具函数（响应封装、ID 生成等）
├── docs/                # 项目文档
├── bin/                 # 编译产物输出目录
├── config.json          # 运行时配置文件（不提交到版本控制）
├── config.example.json  # 配置示例文件（提交到版本控制）
├── go.mod
└── go.sum
```

**规则：**

- `cmd/` 下只放 `main.go`，只做依赖组装和服务启动，**严禁放业务逻辑**。
- 业务逻辑必须在 `internal/` 下实现，禁止其他项目引用。
- 公共工具函数放 `pkg/`，应保持无业务依赖，可独立使用。
- 禁止在项目中创建 `/src` 目录。

---

## 二、命名规范

### 2.1 包名

- 全部小写，无下划线，无驼峰，与目录名一致。
- 简短且具有描述性。

```go
// 正确
package handlers
package services
package utils

// 错误
package handler_utils
package Services
```

### 2.2 文件名

- 全部小写，单词间用下划线分隔，体现其核心职责。
- 命名格式：`{业务模块}_{层级}.go`

```
user_handler.go      // 用户 handler
user_service.go      // 用户 service
user_repository.go   // 用户 repository
```

### 2.3 变量与函数命名

- 使用驼峰命名法（camelCase/PascalCase）。
- 导出标识符（公开）使用大驼峰，包内标识符使用小驼峰。
- 缩写词保持全大写或全小写，如 `userID`、`parseURL`、`HTTPClient`。

```go
// 正确
var userID uint64
func GetUserByID(id uint64) (*models.User, error)
type UserService struct{}

// 错误
var userId uint64
func get_user_by_id(id uint64) (*models.User, error)
```

### 2.4 常量命名

- 使用大驼峰（PascalCase）命名，禁止全大写加下划线风格（除非是真正的全局配置常量）。

```go
const MaxPageSize = 100
const DefaultAvatar = "/assets/default.png"
```

### 2.5 接口命名

- 接口名以功能描述为主，单方法接口可用 `-er` 后缀。

```go
type UserRepository interface { ... }
type Stringer interface { String() string }
```

### 2.6 数据库表名

- 表名使用**单数形式**，全小写，单词间用下划线分隔。
- 通过 `TableName()` 方法显式指定，不依赖 GORM 自动推断。

```go
func (User) TableName() string { return "user" }
func (Poem) TableName() string { return "poem" }
```

---

## 三、代码风格规范

### 3.1 格式化

- 所有代码必须通过 `gofmt` 格式化，提交前必须执行。
- 推荐使用 `goimports` 自动管理 import 分组。

### 3.2 Import 分组

按以下顺序分三组，组间空行隔开：

```go
import (
    // 标准库
    "fmt"
    "net/http"

    // 第三方库
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    // 项目内部包
    "shihai/internal/dto"
    "shihai/internal/services"
    "shihai/pkg/utils"
)
```

### 3.3 注释规范

- 所有导出的类型、函数、方法、变量都必须有注释，以名称开头。
- 复杂逻辑、非直觉代码必须添加行内注释。

```go
// UserService 提供用户相关的业务逻辑处理
type UserService struct { ... }

// GetByID 根据用户 ID 查询用户信息，若用户不存在返回 nil 和错误
func (s *UserService) GetByID(id uint64) (*models.User, error) { ... }
```

### 3.4 函数长度

- 单个函数原则上不超过 **80 行**，超出时应拆分为子函数。
- 函数只做一件事，保持单一职责。

### 3.5 错误处理

- 禁止忽略错误返回值，必须显式处理。
- 不使用 `panic`（除非是程序启动时的不可恢复错误）。

```go
// 正确
user, err := s.repo.FindByID(id)
if err != nil {
    return nil, fmt.Errorf("查询用户失败: %w", err)
}

// 错误
user, _ := s.repo.FindByID(id)
```

---

## 四、分层架构规范

项目严格遵循 **Handler → Service → Repository** 三层架构，各层职责清晰，**禁止跨层调用**。

```
HTTP 请求
    ↓
Handler（参数绑定与校验 → 调用 Service → 返回响应）
    ↓
Service（业务逻辑处理 → 调用 Repository）
    ↓
Repository（数据访问 → 操作数据库）
    ↓
数据库（PostgreSQL via GORM）
```

### 4.1 Handler 层规范

- 只负责：请求参数绑定、参数校验、调用 Service、统一响应返回。
- 不包含任何业务逻辑或 SQL 操作。

```go
func (h *UserHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err.Error())
        return
    }
    user, err := h.userService.Register(&req)
    if err != nil {
        utils.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    utils.Success(c, user)
}
```

### 4.2 Service 层规范

- 包含核心业务逻辑，不直接操作数据库。
- 通过 Repository 接口访问数据，便于测试和替换。
- 跨模块调用通过注入其他 Service 实现，禁止直接调用其他模块的 Repository。

### 4.3 Repository 层规范

- 只负责数据库 CRUD 操作，不包含业务逻辑。
- 方法命名语义化：`FindByID`、`FindAll`、`Create`、`Update`、`DeleteByID`。
- 通过 `*gorm.DB` 执行查询，统一使用 GORM 方法，禁止原始 SQL 字符串拼接（防注入）。

### 4.4 DTO 规范

- 所有 HTTP 请求体和响应体使用 DTO，**禁止直接将 Model 暴露给外部**。
- 请求 DTO 使用 `validator` 标签进行校验约束。

```go
// dto/user_dto.go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6,max=100"`
    Name     string `json:"name" binding:"required,max=50"`
}
```

---

## 五、API 设计规范

### 5.1 URL 设计

- 使用 RESTful 风格，资源名用小写复数名词，单词间用连字符 `-` 分隔。
- 不在 URL 中使用动词。

```
GET    /api/v1/poems            # 获取诗词列表
GET    /api/v1/poems/:id        # 获取单首诗词
POST   /api/v1/poems            # 创建诗词
PUT    /api/v1/poems/:id        # 更新诗词
DELETE /api/v1/poems/:id        # 删除诗词
GET    /api/v1/poems/:id/comments  # 获取诗词评论
```

### 5.2 版本控制

- 所有 API 路径包含版本号前缀 `/api/v1/`。

### 5.3 统一响应格式

所有接口返回统一的 JSON 结构：

```json
{
    "code": 200,
    "message": "success",
    "data": { ... }
}
```

| 字段      | 类型   | 说明               |
|-----------|--------|--------------------|
| `code`    | int    | 业务状态码         |
| `message` | string | 提示信息           |
| `data`    | any    | 响应数据，失败时为 null |

- 通过 `pkg/utils` 中封装的 `utils.Success()`、`utils.Error()`、`utils.BadRequest()` 等方法统一返回。
- **禁止**在 Handler 中直接调用 `c.JSON()` 手动构造响应。

### 5.4 分页参数

分页请求统一使用以下参数：

| 参数       | 类型 | 默认值 | 说明       |
|------------|------|--------|------------|
| `page`     | int  | 1      | 页码（从 1 开始）|
| `pageSize` | int  | 10     | 每页条数，最大 100 |

### 5.5 HTTP 状态码使用

| 状态码 | 场景                           |
|--------|-------------------------------|
| 200    | 请求成功                       |
| 201    | 资源创建成功                   |
| 400    | 请求参数错误                   |
| 401    | 未认证（Token 缺失或无效）       |
| 403    | 无权限                         |
| 404    | 资源不存在                     |
| 500    | 服务器内部错误                 |

---

## 六、错误处理规范

### 6.1 错误包装

使用 `fmt.Errorf` 和 `%w` 动词包装错误，保留调用链：

```go
if err != nil {
    return nil, fmt.Errorf("UserService.Register: 创建用户失败: %w", err)
}
```

### 6.2 业务错误与系统错误区分

- **业务错误**（如用户名已存在）：返回 HTTP 400，附带可读的错误提示。
- **系统错误**（如数据库连接失败）：返回 HTTP 500，日志记录详细信息，响应不暴露内部细节。

### 6.3 禁止 panic

- 禁止在业务代码中使用 `panic`。
- 只在 `main.go` 中初始化失败时允许 `log.Fatal` / `panic`（如数据库连接失败）。

---

## 七、数据模型规范

### 7.1 基础模型

所有业务模型必须内嵌 `BaseModel`，自动获得主键、时间戳和软删除能力：

```go
type BaseModel struct {
    ID        uint64         `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    CreatedBy uint64         `json:"createdBy" gorm:"default:0"`
    UpdatedBy uint64         `json:"updatedBy" gorm:"default:0"`
}
```

### 7.2 主键规范

- 主键统一使用 `uint64` 类型，通过雪花算法自动生成。
- 禁止使用自增整型 ID 暴露给外部。

### 7.3 字段标签规范

每个字段必须同时添加 `json` 和 `gorm` 标签：

```go
Title   string `json:"title" gorm:"not null;size:200;comment:诗词标题"`
Content string `json:"content" gorm:"not null;type:text;comment:诗词内容"`
```

- 密码等敏感字段 JSON 序列化设为 `json:"-"`。
- 关联对象字段添加 `omitempty`，避免空值被序列化。
- `gorm` 标签中必须添加 `comment` 说明字段含义。

### 7.4 软删除

- 统一使用软删除（`gorm.DeletedAt`），**禁止物理删除**业务数据。
- 查询时 GORM 自动过滤软删除记录，无需手动添加条件。

---

## 八、数据库操作规范

### 8.1 使用 GORM

- 统一使用 GORM v2 操作数据库（PostgreSQL）。
- 禁止直接拼接 SQL 字符串，防止 SQL 注入；如需原生 SQL，使用参数化查询：

```go
// 正确
db.Where("username = ?", username).First(&user)

// 错误
db.Where("username = '" + username + "'").First(&user)
```

### 8.2 事务处理

涉及多表写操作时，必须使用事务：

```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&poem).Error; err != nil {
        return err
    }
    if err := tx.Create(&tag).Error; err != nil {
        return err
    }
    return nil
})
```

### 8.3 查询优化

- 按需加载关联关系，使用 `Preload` 而非 N+1 查询。
- 列表查询必须分页，禁止不加 `LIMIT` 的全表查询。
- 高频查询字段确保有数据库索引。

```go
// 正确：分页 + 按需 Preload
db.Preload("Author").Preload("Dynasty").
    Offset(offset).Limit(pageSize).
    Find(&poems)
```

---

## 九、认证与权限规范

### 9.1 JWT 认证

- 使用 JWT（`golang-jwt/jwt/v5`）进行用户身份认证。
- Token 通过 HTTP Header 传递：`Authorization: Bearer <token>`。
- Token 过期时间、密钥等配置从 `config.json` 读取，**严禁硬编码**。

### 9.2 RBAC 权限控制

- 系统采用基于角色的访问控制（RBAC），权限码定义在 `models/permission_codes.go`。
- 通过 `middleware/rbac.go` 中间件拦截权限校验。
- 新增接口时必须在路由注册处明确声明所需权限。

```go
// 路由注册示例
admin := r.Group("/api/v1/admin")
admin.Use(middleware.Auth(), middleware.RequirePermission(models.PermManageUsers))
{
    admin.GET("/users", userHandler.ListUsers)
}
```

### 9.3 密码安全

- 密码存储必须使用 `bcrypt` 加密，**禁止明文存储**。
- 密码字段在 JSON 序列化时必须屏蔽（`json:"-"`）。

---

## 十、日志规范

- 使用标准库 `log` 或项目统一的日志工具，禁止使用 `fmt.Println` 输出运行日志。
- 日志分级：`DEBUG`（开发调试）、`INFO`（关键操作）、`WARN`（异常但可继续）、`ERROR`（需要关注的错误）。
- 日志内容包含：时间、级别、模块、关键信息，**禁止记录用户密码、Token 等敏感信息**。
- 错误日志必须包含足够的上下文，方便排查：

```go
log.Printf("[ERROR] UserService.Login: 用户 %s 登录失败: %v", username, err)
```

---

## 十一、配置管理规范

- 配置文件为 `config.json`，通过 `internal/config/` 包加载并解析为结构体。
- 提供 `config.example.json` 作为配置模板，**`config.json` 不提交到版本控制**（已加入 `.gitignore`）。
- 禁止在代码中硬编码任何配置值（数据库连接串、密钥、端口等）。
- 配置结构体字段必须添加注释说明含义和默认值。

```go
type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    JWT      JWTConfig      `json:"jwt"`
}

type JWTConfig struct {
    Secret     string `json:"secret"`      // JWT 签名密钥
    ExpireHour int    `json:"expireHour"`  // Token 有效期（小时），默认 24
}
```

---

## 十二、Git 提交与分支管理

### 12.1 提交信息

提交信息的格式、类型和示例以 [Git 提交规范](git-commit-guidelines.md) 为准。

后端相关提交的 `scope` 优先使用业务或技术模块名，例如 `backend`、`poem`、`auth`、`rbac`、`database`、`config`。

### 12.2 分支管理

| 分支        | 用途                         |
|-------------|------------------------------|
| `main`      | 生产稳定版本，只接受 PR 合并 |
| `develop`   | 开发集成分支                 |
| `feat/xxx`  | 功能开发分支                 |
| `fix/xxx`   | Bug 修复分支                 |

- 禁止直接向 `main` 分支推送代码。
- 功能开发从 `develop` 拉取分支，完成后通过 PR 合并回 `develop`。

---

## 十三、测试规范

### 13.1 测试文件

- 测试文件与被测文件放同一目录，命名格式：`{文件名}_test.go`。
- 包名使用 `package xxx_test`（黑盒测试）或 `package xxx`（白盒测试）。

### 13.2 测试覆盖

- Service 层核心业务逻辑必须有单元测试。
- Repository 层使用测试数据库或 mock 进行集成测试。
- Handler 层可使用 `httptest` 进行接口测试。

### 13.3 命名规范

```go
func TestUserService_Register_Success(t *testing.T) { ... }
func TestUserService_Register_DuplicateUsername(t *testing.T) { ... }
```

格式：`Test{结构体}_{方法}_{场景}`

### 13.4 运行测试

```bash
# 运行所有测试
go test ./...

# 运行指定包测试并输出覆盖率
go test ./internal/services/... -v -cover
```

---

*最后更新：2026-06-03*
