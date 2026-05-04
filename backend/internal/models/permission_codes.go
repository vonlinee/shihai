package models

// =============================================================================
// 权限编码常量 - 全局唯一来源 (Single Source of Truth)
//
// 设计原则：
//   1. 所有权限编码字符串必须以常量形式定义在本文件中，禁止在其他地方
//      使用形如 "user:create" 的字符串字面量。
//   2. 常量命名规则：Perm + 模块名 + 操作名（驼峰大小写），
//      值规则："module:action"（小写，冒号分隔）。
//   3. 新增 API 接口时，必须执行以下三步：
//      a) 在本文件中追加 PermXxx 常量；
//      b) 在 AllPermissions 中追加对应的 PermissionDef 元数据；
//      c) 在 RolePermissionMap (rbac.go) 中按角色追加授权（admin 自动获取全部，无需追加）。
//   4. 路由注册时使用 RequirePermission(models.PermXxx)，避免硬编码字符串。
//
// 优势：
//   - 编译期检查：拼写错误（如 "usre:list"）直接编译失败；
//   - IDE 重构友好：常量改名/删除可一键定位所有引用；
//   - 易于审计：本文件即权限点全集，便于安全审查与文档生成。
// =============================================================================

// ----------------------------- 用户模块 (user) ------------------------------
const (
	PermUserCreate = "user:create" // 创建新用户
	PermUserRead   = "user:read"   // 查看单个用户信息
	PermUserUpdate = "user:update" // 更新用户信息
	PermUserDelete = "user:delete" // 删除用户
	PermUserList   = "user:list"   // 查看用户列表
)

// ----------------------------- 角色模块 (role) ------------------------------
const (
	PermRoleCreate = "role:create" // 创建新角色
	PermRoleRead   = "role:read"   // 查看单个角色信息
	PermRoleUpdate = "role:update" // 更新角色信息
	PermRoleDelete = "role:delete" // 删除角色
	PermRoleList   = "role:list"   // 查看角色列表
	PermRoleAssign = "role:assign" // 为用户分配角色
)

// -------------------------- 权限模块 (permission) ---------------------------
const (
	PermPermissionList   = "permission:list"   // 查看权限列表
	PermPermissionRead   = "permission:read"   // 查看单个权限详情
	PermPermissionCreate = "permission:create" // 创建权限
	PermPermissionUpdate = "permission:update" // 更新权限信息
	PermPermissionDelete = "permission:delete" // 删除权限
)

// ----------------------------- 诗词模块 (poem) ------------------------------
const (
	PermPoemCreate = "poem:create" // 添加新诗词
	PermPoemRead   = "poem:read"   // 查看诗词详情
	PermPoemUpdate = "poem:update" // 编辑诗词信息
	PermPoemDelete = "poem:delete" // 删除诗词
	PermPoemList   = "poem:list"   // 查看诗词列表
)

// ---------------------------- 评论模块 (comment) ----------------------------
const (
	PermCommentCreate   = "comment:create"   // 发表评论
	PermCommentRead     = "comment:read"     // 查看评论
	PermCommentUpdate   = "comment:update"   // 更新评论
	PermCommentDelete   = "comment:delete"   // 删除评论
	PermCommentList     = "comment:list"     // 查看评论列表
	PermCommentModerate = "comment:moderate" // 审核管理评论
)

// -------------------------- 纠错模块 (correction) ---------------------------
const (
	PermCorrectionCreate = "correction:create" // 提交纠错申请
	PermCorrectionRead   = "correction:read"   // 查看纠错详情
	PermCorrectionUpdate = "correction:update" // 更新纠错信息
	PermCorrectionReview = "correction:review" // 审核纠错申请
	PermCorrectionList   = "correction:list"   // 查看纠错列表
)

// ------------------------ 公告模块 (announcement) ---------------------------
const (
	PermAnnouncementCreate = "announcement:create" // 创建公告
	PermAnnouncementRead   = "announcement:read"   // 查看公告
	PermAnnouncementUpdate = "announcement:update" // 更新公告
	PermAnnouncementDelete = "announcement:delete" // 删除公告
	PermAnnouncementList   = "announcement:list"   // 查看公告列表
)

// ---------------------------- 论坛模块 (forum) ------------------------------
const (
	PermForumCreate   = "forum:create"   // 发表论坛帖子
	PermForumRead     = "forum:read"     // 查看论坛帖子
	PermForumUpdate   = "forum:update"   // 更新论坛帖子
	PermForumDelete   = "forum:delete"   // 删除论坛帖子
	PermForumList     = "forum:list"     // 查看帖子列表
	PermForumModerate = "forum:moderate" // 审核管理论坛帖子
)

// ----------------------------- 题目模块 (quiz) ------------------------------
const (
	PermQuizCreate = "quiz:create" // 创建问答题目
	PermQuizRead   = "quiz:read"   // 查看问答题目
	PermQuizUpdate = "quiz:update" // 更新问答题目
	PermQuizDelete = "quiz:delete" // 删除问答题目
	PermQuizList   = "quiz:list"   // 查看题目列表
)

// --------------------------- 反馈模块 (feedback) ----------------------------
const (
	PermFeedbackCreate = "feedback:create" // 提交意见反馈
	PermFeedbackRead   = "feedback:read"   // 查看反馈详情
	PermFeedbackUpdate = "feedback:update" // 更新反馈状态
	PermFeedbackDelete = "feedback:delete" // 删除反馈
	PermFeedbackList   = "feedback:list"   // 查看反馈列表
)

// ---------------------------- 系统模块 (system) -----------------------------
const (
	PermSystemConfig = "system:config" // 系统配置管理
	PermSystemLog    = "system:log"    // 查看系统日志
	PermSystemBackup = "system:backup" // 系统备份管理
)

// ----------------------------- 模块名称常量 ---------------------------------
// 用于 PermissionDef.Module 字段，避免在元数据中再次出现字符串字面量
const (
	ModuleUser         = "user"         // 用户模块
	ModuleRole         = "role"         // 角色模块
	ModulePermission   = "permission"   // 权限模块
	ModulePoem         = "poem"         // 诗词模块
	ModuleComment      = "comment"      // 评论模块
	ModuleCorrection   = "correction"   // 纠错模块
	ModuleAnnouncement = "announcement" // 公告模块
	ModuleForum        = "forum"        // 论坛模块
	ModuleQuiz         = "quiz"         // 题目模块
	ModuleFeedback     = "feedback"     // 反馈模块
	ModuleSystem       = "system"       // 系统模块
)

// PermissionDef 权限元数据定义
// 用于系统启动时将权限点写入 permission 表，包含展示所需的中文名称、
// 所属模块和详细描述。Code 字段必须严格使用本文件中定义的 PermXxx 常量。
type PermissionDef struct {
	Code        string // 权限编码，必须使用 PermXxx 常量，禁止字面量
	Name        string // 权限中文显示名称（用于管理后台展示）
	Module      string // 所属模块，必须使用 ModuleXxx 常量
	Description string // 权限详细描述，用于鼠标悬浮提示和文档生成
}

// AllPermissions 全部权限定义（系统启动时自动同步到 permission 表）
//
// 维护规则：
//   - 新增权限：先在上方常量区追加 PermXxx 常量，再在此切片追加对应元数据；
//   - 删除权限：同时删除常量定义、本切片中的元数据，并清理 RolePermissionMap 中的引用；
//   - 修改权限：修改 PermissionDef 中的 Name/Module/Description 即可，
//     重启服务后会通过 InitDefaultRolesAndPermissions 自动 UPSERT 到数据库。
//
// admin 角色启动时会自动获取本切片中的所有权限，无需在 RolePermissionMap 中重复声明。
var AllPermissions = []PermissionDef{
	// user
	{PermUserCreate, "创建用户", ModuleUser, "创建新用户"},
	{PermUserRead, "查看用户", ModuleUser, "查看用户信息"},
	{PermUserUpdate, "更新用户", ModuleUser, "更新用户信息"},
	{PermUserDelete, "删除用户", ModuleUser, "删除用户"},
	{PermUserList, "用户列表", ModuleUser, "查看用户列表"},
	// role
	{PermRoleCreate, "创建角色", ModuleRole, "创建新角色"},
	{PermRoleRead, "查看角色", ModuleRole, "查看角色信息"},
	{PermRoleUpdate, "更新角色", ModuleRole, "更新角色信息"},
	{PermRoleDelete, "删除角色", ModuleRole, "删除角色"},
	{PermRoleList, "角色列表", ModuleRole, "查看角色列表"},
	{PermRoleAssign, "分配角色", ModuleRole, "为用户分配角色"},
	// permission
	{PermPermissionList, "权限列表", ModulePermission, "查看权限列表"},
	{PermPermissionRead, "查看权限", ModulePermission, "查看权限详情"},
	{PermPermissionCreate, "创建权限", ModulePermission, "创建新权限"},
	{PermPermissionUpdate, "更新权限", ModulePermission, "更新权限信息"},
	{PermPermissionDelete, "删除权限", ModulePermission, "删除权限"},
	// poem
	{PermPoemCreate, "创建诗词", ModulePoem, "添加新诗词"},
	{PermPoemRead, "查看诗词", ModulePoem, "查看诗词详情"},
	{PermPoemUpdate, "更新诗词", ModulePoem, "编辑诗词信息"},
	{PermPoemDelete, "删除诗词", ModulePoem, "删除诗词"},
	{PermPoemList, "诗词列表", ModulePoem, "查看诗词列表"},
	// comment
	{PermCommentCreate, "创建评论", ModuleComment, "发表评论"},
	{PermCommentRead, "查看评论", ModuleComment, "查看评论"},
	{PermCommentUpdate, "更新评论", ModuleComment, "更新评论"},
	{PermCommentDelete, "删除评论", ModuleComment, "删除评论"},
	{PermCommentList, "评论列表", ModuleComment, "查看评论列表"},
	{PermCommentModerate, "审核评论", ModuleComment, "审核管理评论"},
	// correction
	{PermCorrectionCreate, "提交纠错", ModuleCorrection, "提交纠错申请"},
	{PermCorrectionRead, "查看纠错", ModuleCorrection, "查看纠错详情"},
	{PermCorrectionUpdate, "更新纠错", ModuleCorrection, "更新纠错信息"},
	{PermCorrectionReview, "审核纠错", ModuleCorrection, "审核纠错申请"},
	{PermCorrectionList, "纠错列表", ModuleCorrection, "查看纠错列表"},
	// announcement
	{PermAnnouncementCreate, "创建公告", ModuleAnnouncement, "创建公告"},
	{PermAnnouncementRead, "查看公告", ModuleAnnouncement, "查看公告"},
	{PermAnnouncementUpdate, "更新公告", ModuleAnnouncement, "更新公告"},
	{PermAnnouncementDelete, "删除公告", ModuleAnnouncement, "删除公告"},
	{PermAnnouncementList, "公告列表", ModuleAnnouncement, "查看公告列表"},
	// forum
	{PermForumCreate, "发表帖子", ModuleForum, "发表论坛帖子"},
	{PermForumRead, "查看帖子", ModuleForum, "查看论坛帖子"},
	{PermForumUpdate, "更新帖子", ModuleForum, "更新论坛帖子"},
	{PermForumDelete, "删除帖子", ModuleForum, "删除论坛帖子"},
	{PermForumList, "帖子列表", ModuleForum, "查看帖子列表"},
	{PermForumModerate, "审核帖子", ModuleForum, "审核管理论坛帖子"},
	// quiz
	{PermQuizCreate, "创建题目", ModuleQuiz, "创建问答题目"},
	{PermQuizRead, "查看题目", ModuleQuiz, "查看问答题目"},
	{PermQuizUpdate, "更新题目", ModuleQuiz, "更新问答题目"},
	{PermQuizDelete, "删除题目", ModuleQuiz, "删除问答题目"},
	{PermQuizList, "题目列表", ModuleQuiz, "查看题目列表"},
	// feedback
	{PermFeedbackCreate, "提交反馈", ModuleFeedback, "提交意见反馈"},
	{PermFeedbackRead, "查看反馈", ModuleFeedback, "查看反馈详情"},
	{PermFeedbackUpdate, "更新反馈", ModuleFeedback, "更新反馈状态"},
	{PermFeedbackDelete, "删除反馈", ModuleFeedback, "删除反馈"},
	{PermFeedbackList, "反馈列表", ModuleFeedback, "查看反馈列表"},
	// system
	{PermSystemConfig, "系统配置", ModuleSystem, "系统配置管理"},
	{PermSystemLog, "系统日志", ModuleSystem, "查看系统日志"},
	{PermSystemBackup, "系统备份", ModuleSystem, "系统备份管理"},
}
