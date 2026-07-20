# 论坛完整链路设计

## 目标

补齐论坛发帖、回帖、前台浏览和后台管理的完整业务链路，使论坛从“部分实现”进入当前可用范围。

## 后端

### 分层

新增标准后端分层：

- `ForumRepository`：封装帖子和回复的持久化查询、写入、业务删除、置顶、浏览量和回复数维护。
- `ForumService`：负责输入裁剪、作者或审核权限判断、父回复归属校验、分页边界和 DTO 映射。
- `ForumHandler`：负责 Gin 参数绑定、当前用户读取、状态码映射和统一响应封装。

### 公开 API

公开读接口挂载在 `/api/forum`，不要求登录：

- `GET /api/forum/posts`：分页查询未删除帖子。查询参数：`page`、`pageSize`、`keyword`。
- `GET /api/forum/posts/:id`：查询未删除帖子详情，并增加浏览量。
- `GET /api/forum/posts/:id/replies`：分页查询未删除回复。查询参数：`page`、`pageSize`。

### 登录用户 API

登录写接口挂载在 `/api/forum`，要求 `Authorization: Bearer <token>`：

- `POST /api/forum/posts`：发表帖子，要求 `forum:create` 权限。请求体：`title`、`content`。
- `PUT /api/forum/posts/:id`：更新自己的帖子。请求体：`title`、`content`，字段为空时保持原值。
- `DELETE /api/forum/posts/:id`：删除自己的帖子。
- `POST /api/forum/posts/:id/replies`：回复帖子，要求 `forum:create` 权限。请求体：`content`、可选 `parentId`。
- `DELETE /api/forum/replies/:id`：删除自己的回复。

### 后台 API

后台论坛管理接口挂载在 `/api/admin/forum`，要求登录和 `forum:moderate` 权限，不套用 `/api/admin` 分组的 admin/editor 角色门槛：

- `GET /api/admin/forum/posts`：分页查询帖子管理列表。查询参数：`page`、`pageSize`、`keyword`、`includeDeleted`。
- `PUT /api/admin/forum/posts/:id/pin`：设置置顶状态。请求体：`isPinned`。
- `DELETE /api/admin/forum/posts/:id`：审核删除帖子。
- `DELETE /api/admin/forum/replies/:id`：审核删除回复。

### DTO

请求 DTO：

- `ForumPostCreateRequest`：`title` 2-200 个字符，`content` 2-10000 个字符。
- `ForumPostUpdateRequest`：可选 `title`、`content`。
- `ForumReplyCreateRequest`：`content` 1-5000 个字符，可选 `parentId` 使用 `RequestID` 解析。
- `ForumPostPinRequest`：`isPinned` 布尔值。

响应 DTO：

- `ForumPostResponse`：返回帖子 ID、用户 ID、用户摘要、标题、正文、浏览量、回复数、置顶状态、删除状态和时间字段。
- `ForumReplyResponse`：返回回复 ID、帖子 ID、用户 ID、用户摘要、正文、父回复 ID、父回复摘要、删除状态和时间字段。
- 所有 Snowflake ID 对外按字符串序列化，避免前端精度丢失。

### 业务规则

- 前台列表、详情和回复列表只返回未删除内容。
- 后台列表默认不包含已删除帖子，传 `includeDeleted=true` 时包含。
- 普通用户只能更新或删除自己的帖子和回复。
- 审核删除和置顶由后台路由权限控制。
- 回复必须属于有效且未删除的帖子。
- 楼中楼父回复必须存在、未删除，并且属于同一帖子。
- 创建回复与帖子回复数递增在同一事务内完成。
- 删除回复与帖子回复数递减在同一事务内完成，回复数不会递减到负数。

### 错误映射

- 路径参数格式错误返回 400。
- 未登录返回 401。
- 非作者执行用户自助修改或删除返回 403。
- 帖子或回复不存在返回 404。
- 父回复归属不合法、正文为空等业务校验失败返回 422。
- 未知持久化错误返回 500，客户端不暴露数据库细节。

## 前端

### 前台页面

`/forum` 展示论坛列表：

- 关键词搜索。
- 置顶帖突出展示。
- 帖子列表显示回复数、浏览量、标题摘要、作者和更新时间。
- 登录用户可以在侧栏发表主题。

`/forum/:id` 展示帖子详情：

- 展示帖子正文、作者、浏览量和回复数。
- 登录用户可以回复帖子或楼中楼回复。
- 作者可以删除自己的帖子和回复。
- 删除操作使用项目确认弹窗。

### 后台页面

`/admin/forum` 展示论坛管理列表：

- 关键词搜索。
- 可切换是否包含已删除帖子。
- 表格展示帖子、作者、状态、回复/浏览和更新时间。
- 支持查看帖子、置顶或取消置顶、删除帖子。

### 前端数据层

新增前端论坛数据层：

- `frontend/src/services/forumService.ts`：封装前台论坛 API。
- `frontend/src/hooks/useForum.ts`：封装 TanStack Query 查询和写操作。
- `frontend/src/pages/ForumPage.tsx`：论坛列表和详情页面。
- `frontend/src/pages/admin/ForumPage.tsx`：后台论坛管理页面。

扩展后台数据层：

- `adminService.getForumPosts`
- `adminService.setForumPostPinned`
- `adminService.deleteForumPost`
- `useAdminForumPosts`
- `useAdminSetForumPostPinned`
- `useAdminDeleteForumPost`

扩展共享类型：

- `ForumUser`
- `ForumPost`
- `ForumReply`

## 文档

- Swagger 注解在 `backend/internal/apidocs/annotations_public.go` 和 `backend/internal/apidocs/annotations_admin.go` 维护。
- Swagger 生成文件位于 `backend/docs/swagger/`。
- 长期需求基线和系统设计同步将论坛移动到当前能力范围。

## 验证

后端验证：

- `gofmt` 格式化新增和修改的 Go 文件。
- `go test ./...` 通过。

前端验证：

- `npm run build` 通过 TypeScript 和 Vite 构建。
