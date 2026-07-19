package apidocs

import (
	"shihai/internal/dto"
	"shihai/pkg/utils"
)

var (
	_ dto.AdminCreateUserRequest
	_ dto.CorrectionResponse
	_ utils.Response
)

// swaggerAdminListUsers documents GET /api/admin/users.
// @Summary 查询用户列表
// @Tags 后台用户
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "关键词"
// @Param role query string false "角色"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/users [get]
func swaggerAdminListUsers() {}

// swaggerAdminGetUser documents GET /api/admin/users/{id}.
// @Summary 获取用户详情
// @Tags 后台用户
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/admin/users/{id} [get]
func swaggerAdminGetUser() {}

// swaggerAdminCreateUser documents POST /api/admin/users.
// @Summary 创建用户
// @Tags 后台用户
// @Security BearerAuth
// @Param request body dto.AdminCreateUserRequest true "用户信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/users [post]
func swaggerAdminCreateUser() {}

// swaggerAdminDeleteUser documents DELETE /api/admin/users/{id}.
// @Summary 删除用户
// @Tags 后台用户
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/users/{id} [delete]
func swaggerAdminDeleteUser() {}

// swaggerAdminCreatePoem documents POST /api/admin/poems.
// @Summary 创建诗词
// @Tags 后台诗词
// @Security BearerAuth
// @Param request body dto.PoemCreateRequest true "诗词信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems [post]
func swaggerAdminCreatePoem() {}

// swaggerAdminUpdatePoem documents PUT /api/admin/poems/{id}.
// @Summary 更新诗词
// @Tags 后台诗词
// @Security BearerAuth
// @Param id path string true "诗词 ID"
// @Param request body dto.PoemUpdateRequest true "诗词信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems/{id} [put]
func swaggerAdminUpdatePoem() {}

// swaggerAdminBatchDeletePoems documents DELETE /api/admin/poems.
// @Summary 批量删除诗词
// @Tags 后台诗词
// @Security BearerAuth
// @Param request body dto.BatchDeleteRequest true "诗词 ID 列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems [delete]
func swaggerAdminBatchDeletePoems() {}

// swaggerAdminDeletePoem documents DELETE /api/admin/poems/{id}.
// @Summary 删除诗词
// @Tags 后台诗词
// @Security BearerAuth
// @Param id path string true "诗词 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems/{id} [delete]
func swaggerAdminDeletePoem() {}

// swaggerAdminListPoemAnnotations documents GET /api/admin/poems/{id}/annotations.
// @Summary 查询诗词标注
// @Tags 后台诗词标注
// @Security BearerAuth
// @Param id path string true "诗词 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems/{id}/annotations [get]
func swaggerAdminListPoemAnnotations() {}

// swaggerAdminCreatePoemAnnotation documents POST /api/admin/poems/{id}/annotations.
// @Summary 创建诗词标注
// @Tags 后台诗词标注
// @Security BearerAuth
// @Param id path string true "诗词 ID"
// @Param request body dto.PoemAnnotationCreateRequest true "标注信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poems/{id}/annotations [post]
func swaggerAdminCreatePoemAnnotation() {}

// swaggerAdminUpdatePoemAnnotation documents PUT /api/admin/poem-annotations/{id}.
// @Summary 更新诗词标注
// @Tags 后台诗词标注
// @Security BearerAuth
// @Param id path string true "标注 ID"
// @Param request body dto.PoemAnnotationUpdateRequest true "标注信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poem-annotations/{id} [put]
func swaggerAdminUpdatePoemAnnotation() {}

// swaggerAdminDeletePoemAnnotation documents DELETE /api/admin/poem-annotations/{id}.
// @Summary 删除诗词标注
// @Tags 后台诗词标注
// @Security BearerAuth
// @Param id path string true "标注 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poem-annotations/{id} [delete]
func swaggerAdminDeletePoemAnnotation() {}

// swaggerAdminConvertText documents POST /api/admin/text-conversion.
// @Summary 批量简繁转换
// @Tags 后台文本转换
// @Security BearerAuth
// @Param request body dto.TextConversionRequest true "转换信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/text-conversion [post]
func swaggerAdminConvertText() {}

// swaggerAdminCreateDynasty documents POST /api/admin/dynasties.
// @Summary 创建朝代
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.DynastyCreateRequest true "朝代信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/dynasties [post]
func swaggerAdminCreateDynasty() {}

// swaggerAdminUpdateDynasty documents PUT /api/admin/dynasties/{id}.
// @Summary 更新朝代
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "朝代 ID"
// @Param request body dto.DynastyUpdateRequest true "朝代信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/dynasties/{id} [put]
func swaggerAdminUpdateDynasty() {}

// swaggerAdminBatchDeleteDynasties documents DELETE /api/admin/dynasties.
// @Summary 批量删除朝代
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.BatchDeleteRequest true "朝代 ID 列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/dynasties [delete]
func swaggerAdminBatchDeleteDynasties() {}

// swaggerAdminDeleteDynasty documents DELETE /api/admin/dynasties/{id}.
// @Summary 删除朝代
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "朝代 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/dynasties/{id} [delete]
func swaggerAdminDeleteDynasty() {}

// swaggerAdminCreatePoet documents POST /api/admin/poets.
// @Summary 创建诗人
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.PoetCreateRequest true "诗人信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poets [post]
func swaggerAdminCreatePoet() {}

// swaggerAdminUpdatePoet documents PUT /api/admin/poets/{id}.
// @Summary 更新诗人
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "诗人 ID"
// @Param request body dto.PoetUpdateRequest true "诗人信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poets/{id} [put]
func swaggerAdminUpdatePoet() {}

// swaggerAdminBatchDeletePoets documents DELETE /api/admin/poets.
// @Summary 批量删除诗人
// @Tags 后台基础数据
// @Security BearerAuth
// @Param request body dto.BatchDeleteRequest true "诗人 ID 列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poets [delete]
func swaggerAdminBatchDeletePoets() {}

// swaggerAdminDeletePoet documents DELETE /api/admin/poets/{id}.
// @Summary 删除诗人
// @Tags 后台基础数据
// @Security BearerAuth
// @Param id path string true "诗人 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/poets/{id} [delete]
func swaggerAdminDeletePoet() {}

// swaggerAdminListCollections documents GET /api/admin/work-collections.
// @Summary 查询作品集列表
// @Tags 后台作品集
// @Security BearerAuth
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键词"
// @Param published query bool false "发布状态"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections [get]
func swaggerAdminListCollections() {}

// swaggerAdminGetCollection documents GET /api/admin/work-collections/{id}.
// @Summary 获取作品集详情
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/admin/work-collections/{id} [get]
func swaggerAdminGetCollection() {}

// swaggerAdminCreateCollection documents POST /api/admin/work-collections.
// @Summary 创建作品集
// @Tags 后台作品集
// @Security BearerAuth
// @Param request body dto.WorkCollectionCreateRequest true "作品集信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections [post]
func swaggerAdminCreateCollection() {}

// swaggerAdminUpdateCollection documents PUT /api/admin/work-collections/{id}.
// @Summary 更新作品集
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Param request body dto.WorkCollectionUpdateRequest true "作品集信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id} [put]
func swaggerAdminUpdateCollection() {}

// swaggerAdminDeleteCollection documents DELETE /api/admin/work-collections/{id}.
// @Summary 删除作品集
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id} [delete]
func swaggerAdminDeleteCollection() {}

// swaggerAdminAddCollectionItem documents POST /api/admin/work-collections/{id}/items.
// @Summary 添加作品集条目
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Param request body dto.WorkCollectionItemCreateRequest true "条目信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id}/items [post]
func swaggerAdminAddCollectionItem() {}

// swaggerAdminUpdateCollectionItem documents PUT /api/admin/work-collections/{id}/items/{itemId}.
// @Summary 更新作品集条目
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Param itemId path string true "条目 ID"
// @Param request body dto.WorkCollectionItemUpdateRequest true "条目信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id}/items/{itemId} [put]
func swaggerAdminUpdateCollectionItem() {}

// swaggerAdminDeleteCollectionItem documents DELETE /api/admin/work-collections/{id}/items/{itemId}.
// @Summary 删除作品集条目
// @Tags 后台作品集
// @Security BearerAuth
// @Param id path string true "作品集 ID"
// @Param itemId path string true "条目 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/work-collections/{id}/items/{itemId} [delete]
func swaggerAdminDeleteCollectionItem() {}

// swaggerAdminCreateAnnouncement documents POST /api/admin/announcements.
// @Summary 创建公告
// @Tags 后台公告
// @Security BearerAuth
// @Param request body object true "公告信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/announcements [post]
func swaggerAdminCreateAnnouncement() {}

// swaggerAdminUpdateAnnouncement documents PUT /api/admin/announcements/{id}.
// @Summary 更新公告
// @Tags 后台公告
// @Security BearerAuth
// @Param id path string true "公告 ID"
// @Param request body object true "公告信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/announcements/{id} [put]
func swaggerAdminUpdateAnnouncement() {}

// swaggerAdminDeleteAnnouncement documents DELETE /api/admin/announcements/{id}.
// @Summary 删除公告
// @Tags 后台公告
// @Security BearerAuth
// @Param id path string true "公告 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/announcements/{id} [delete]
func swaggerAdminDeleteAnnouncement() {}

// swaggerAdminListComments documents GET /api/admin/comments/all.
// @Summary 查询全部评论
// @Tags 后台评论
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/comments/all [get]
func swaggerAdminListComments() {}

// swaggerAdminListCorrections documents GET /api/admin/corrections.
// @Summary 查询纠错列表
// @Tags 后台纠错
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "关键词"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/admin/corrections [get]
func swaggerAdminListCorrections() {}
