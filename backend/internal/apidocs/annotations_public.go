package apidocs

import (
	"shihai/internal/dto"
	"shihai/pkg/utils"
)

var (
	_ dto.RegisterRequest
	_ dto.ForumPostCreateRequest
	_ dto.ForumPostUpdateRequest
	_ dto.ForumReplyCreateRequest
	_ utils.Response
)

// swaggerHealth documents GET /health.
// @Summary 服务健康检查
// @Tags 系统
// @Success 200 {object} map[string]string
// @Router /health [get]
func swaggerHealth() {}

// swaggerRegister documents POST /api/auth/register.
// @Summary 用户注册
// @Tags 认证
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/auth/register [post]
func swaggerRegister() {}

// swaggerLogin documents POST /api/auth/login.
// @Summary 用户登录
// @Tags 认证
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/auth/login [post]
func swaggerLogin() {}

// swaggerListPoems documents GET /api/poems.
// @Summary 查询诗词列表
// @Tags 诗词
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "关键词"
// @Param dynasty query string false "朝代"
// @Param author query string false "作者"
// @Param genre query string false "体裁"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/poems [get]
func swaggerListPoems() {}

// swaggerRandomPoems documents GET /api/poems/random.
// @Summary 随机获取诗词
// @Tags 诗词
// @Param limit query int false "返回数量" default(5)
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/poems/random [get]
func swaggerRandomPoems() {}

// swaggerGetPoem documents GET /api/poems/{id}.
// @Summary 获取诗词详情
// @Tags 诗词
// @Param id path string true "诗词 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/poems/{id} [get]
func swaggerGetPoem() {}

// swaggerLikePoem documents POST /api/poems/{id}/like.
// @Summary 点赞诗词
// @Tags 诗词
// @Param id path string true "诗词 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/poems/{id}/like [post]
func swaggerLikePoem() {}

// swaggerListDynasties documents GET /api/dynasties.
// @Summary 查询朝代列表
// @Tags 诗词基础数据
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/dynasties [get]
func swaggerListDynasties() {}

// swaggerListPoets documents GET /api/poets.
// @Summary 查询诗人列表
// @Tags 诗词基础数据
// @Param keyword query string false "关键词"
// @Param dynastyId query string false "朝代 ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/poets [get]
func swaggerListPoets() {}

// swaggerListGenres documents GET /api/genres.
// @Summary 查询体裁列表
// @Tags 诗词基础数据
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/genres [get]
func swaggerListGenres() {}

// swaggerListAnnouncements documents GET /api/announcements.
// @Summary 查询公告列表
// @Tags 公告
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param pinned query bool false "仅置顶公告"
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/announcements [get]
func swaggerListAnnouncements() {}

// swaggerGetAnnouncement documents GET /api/announcements/{id}.
// @Summary 获取公告详情
// @Tags 公告
// @Param id path string true "公告 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/announcements/{id} [get]
func swaggerGetAnnouncement() {}

// swaggerListComments documents GET /api/comments.
// @Summary 查询诗词评论
// @Tags 评论
// @Param poemId query string true "诗词 ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments [get]
func swaggerListComments() {}

// swaggerVoteComment documents POST /api/comments/vote.
// @Summary 评论投票
// @Tags 评论
// @Param X-Visitor-ID header string false "游客标识"
// @Param request body dto.CommentVoteRequest true "投票信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/comments/vote [post]
func swaggerVoteComment() {}

// swaggerGetProfile documents GET /api/user/profile.
// @Summary 获取当前用户资料
// @Tags 用户
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/user/profile [get]
func swaggerGetProfile() {}

// swaggerUpdateProfile documents PUT /api/user/profile.
// @Summary 更新当前用户资料
// @Tags 用户
// @Security BearerAuth
// @Param request body dto.UpdateUserRequest true "用户资料"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/user/profile [put]
func swaggerUpdateProfile() {}

// swaggerChangePassword documents PUT /api/user/password.
// @Summary 修改当前用户密码
// @Tags 用户
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "密码信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/user/password [put]
func swaggerChangePassword() {}

// swaggerCreateComment documents POST /api/comments.
// @Summary 创建评论
// @Tags 评论
// @Security BearerAuth
// @Param request body dto.CommentCreateRequest true "评论信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/comments [post]
func swaggerCreateComment() {}

// swaggerDeleteComment documents DELETE /api/comments/{id}.
// @Summary 删除评论
// @Tags 评论
// @Security BearerAuth
// @Param id path string true "评论 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/comments/{id} [delete]
func swaggerDeleteComment() {}

// swaggerListForumPosts documents GET /api/forum/posts.
// @Summary 查询论坛帖子列表
// @Tags 论坛
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param keyword query string false "关键词"
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts [get]
func swaggerListForumPosts() {}

// swaggerGetForumPost documents GET /api/forum/posts/{id}.
// @Summary 获取论坛帖子详情
// @Tags 论坛
// @Param id path string true "帖子 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/posts/{id} [get]
func swaggerGetForumPost() {}

// swaggerListForumReplies documents GET /api/forum/posts/{id}/replies.
// @Summary 查询论坛帖子回复
// @Tags 论坛
// @Param id path string true "帖子 ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(50)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/posts/{id}/replies [get]
func swaggerListForumReplies() {}

// swaggerCreateForumPost documents POST /api/forum/posts.
// @Summary 创建论坛帖子
// @Tags 论坛
// @Security BearerAuth
// @Param request body dto.ForumPostCreateRequest true "帖子信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /api/forum/posts [post]
func swaggerCreateForumPost() {}

// swaggerUpdateForumPost documents PUT /api/forum/posts/{id}.
// @Summary 更新自己的论坛帖子
// @Tags 论坛
// @Security BearerAuth
// @Param id path string true "帖子 ID"
// @Param request body dto.ForumPostUpdateRequest true "帖子信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/posts/{id} [put]
func swaggerUpdateForumPost() {}

// swaggerDeleteForumPost documents DELETE /api/forum/posts/{id}.
// @Summary 删除自己的论坛帖子
// @Tags 论坛
// @Security BearerAuth
// @Param id path string true "帖子 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/posts/{id} [delete]
func swaggerDeleteForumPost() {}

// swaggerCreateForumReply documents POST /api/forum/posts/{id}/replies.
// @Summary 创建论坛回复
// @Tags 论坛
// @Security BearerAuth
// @Param id path string true "帖子 ID"
// @Param request body dto.ForumReplyCreateRequest true "回复信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /api/forum/posts/{id}/replies [post]
func swaggerCreateForumReply() {}

// swaggerDeleteForumReply documents DELETE /api/forum/replies/{id}.
// @Summary 删除自己的论坛回复
// @Tags 论坛
// @Security BearerAuth
// @Param id path string true "回复 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/forum/replies/{id} [delete]
func swaggerDeleteForumReply() {}
