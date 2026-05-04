package dto

import "time"

// PoemListRequest 诗词列表请求，支持分页、关键词搜索和多种筛选条件
type PoemListRequest struct {
	Page     int    `form:"page,default=1"`      // 页码，默认第1页
	PageSize int    `form:"pageSize,default=10"` // 每页数量，默认10条
	Keyword  string `form:"keyword"`             // 搜索关键词，匹配标题或内容
	Dynasty  string `form:"dynasty"`             // 按朝代筛选
	Author   string `form:"author"`              // 按作者筛选
	Genre    string `form:"genre"`               // 按体裁筛选
}

// PoemCreateRequest 创建诗词请求
// 支持两种关联方式：通过 ID（AuthorID/DynastyID）或通过名称（AuthorName/DynastyName）
// 当 ID 为 0 且 Name 非空时，系统会自动查找或创建对应的诗人/朝代
type PoemCreateRequest struct {
	Title        string `json:"title" binding:"required,max=200"` // 诗词标题，必填，最多200字符
	Content      string `json:"content" binding:"required"`       // 诗词正文，必填
	AuthorID     uint64 `json:"authorId"`                         // 作者ID，与 AuthorName 二选一
	AuthorName   string `json:"authorName"`                       // 作者姓名，与 AuthorID 二选一，支持自动创建
	DynastyID    uint64 `json:"dynastyId"`                        // 朝代ID，与 DynastyName 二选一
	DynastyName  string `json:"dynastyName"`                      // 朝代名称，与 DynastyID 二选一，支持自动创建
	Genre        string `json:"genre" binding:"max=50"`           // 体裁，如"五言绝句"、"七言律诗"
	Translation  string `json:"translation"`                      // 现代汉语译文
	Appreciation string `json:"appreciation"`                     // 诗词赏析与解读
	Annotation   string `json:"annotation"`                       // 字词注释说明
	AudioURL     string `json:"audioUrl" binding:"max=500"`       // 朗读音频地址
	CoverImage   string `json:"coverImage" binding:"max=500"`     // 封面图片地址
}

// PoemUpdateRequest 更新诗词请求，所有字段均为可选，仅更新提供的字段
type PoemUpdateRequest struct {
	Title        string `json:"title" binding:"max=200"`      // 诗词标题
	Content      string `json:"content"`                      // 诗词正文
	AuthorID     uint64 `json:"authorId"`                     // 作者ID
	DynastyID    uint64 `json:"dynastyId"`                    // 朝代ID
	Genre        string `json:"genre" binding:"max=50"`       // 体裁
	Translation  string `json:"translation"`                  // 译文
	Appreciation string `json:"appreciation"`                 // 赏析
	Annotation   string `json:"annotation"`                   // 注释
	AudioURL     string `json:"audioUrl" binding:"max=500"`   // 音频地址
	CoverImage   string `json:"coverImage" binding:"max=500"` // 封面图地址
}

// PoemResponse 诗词响应，包含完整的诗词信息及关联的作者和朝代
type PoemResponse struct {
	ID           uint64          `json:"id"`                // 诗词ID
	Title        string          `json:"title"`             // 诗词标题
	Content      string          `json:"content"`           // 诗词正文
	AuthorID     uint64          `json:"authorId"`          // 作者ID
	Author       PoetResponse    `json:"author,omitempty"`  // 作者信息，预加载时返回
	DynastyID    uint64          `json:"dynastyId"`         // 朝代ID
	Dynasty      DynastyResponse `json:"dynasty,omitempty"` // 朝代信息，预加载时返回
	Genre        string          `json:"genre"`             // 体裁
	Translation  string          `json:"translation"`       // 译文
	Appreciation string          `json:"appreciation"`      // 赏析
	Annotation   string          `json:"annotation"`        // 注释
	AudioURL     string          `json:"audioUrl"`          // 音频地址
	CoverImage   string          `json:"coverImage"`        // 封面图地址
	Views        int             `json:"views"`             // 浏览量
	Likes        int             `json:"likes"`             // 点赞数
	Dislikes     int             `json:"dislikes"`          // 点踩数
	Favorites    int             `json:"favorites"`         // 收藏数
	CreatedAt    time.Time       `json:"createdAt"`         // 创建时间
	UpdatedAt    time.Time       `json:"updatedAt"`         // 更新时间
}

// PoetResponse 诗人响应，包含诗人基本信息
type PoetResponse struct {
	ID        uint64 `json:"id"`        // 诗人ID
	Name      string `json:"name"`      // 诗人姓名
	DynastyID uint64 `json:"dynastyId"` // 所属朝代ID
	Biography string `json:"biography"` // 生平简介
	Avatar    string `json:"avatar"`    // 头像URL
}

// DynastyResponse 朝代响应，包含朝代基本信息
type DynastyResponse struct {
	ID          uint64 `json:"id"`          // 朝代ID
	Name        string `json:"name"`        // 朝代名称
	Period      string `json:"period"`      // 时期描述
	Description string `json:"description"` // 朝代描述
}

// DynastyCreateRequest 创建朝代请求
type DynastyCreateRequest struct {
	Name        string `json:"name" binding:"required,max=50"` // 朝代名称，必填
	Period      string `json:"period" binding:"max=100"`       // 时期描述，如"618-907"
	Description string `json:"description"`                    // 朝代详细描述
}

// DynastyUpdateRequest 更新朝代请求，所有字段均为可选
type DynastyUpdateRequest struct {
	Name        string `json:"name" binding:"max=50"`    // 朝代名称
	Period      string `json:"period" binding:"max=100"` // 时期描述
	Description string `json:"description"`              // 朝代描述
}

// PoetCreateRequest 创建诗人请求
type PoetCreateRequest struct {
	Name      string `json:"name" binding:"required,max=50"` // 诗人姓名，必填
	DynastyID uint64 `json:"dynastyId"`                      // 所属朝代ID
	Biography string `json:"biography"`                      // 生平简介
	Avatar    string `json:"avatar" binding:"max=255"`       // 头像URL
	BirthYear int    `json:"birthYear"`                      // 出生年份
	DeathYear int    `json:"deathYear"`                      // 逝世年份
}

// PoetUpdateRequest 更新诗人请求，所有字段均为可选
type PoetUpdateRequest struct {
	Name      string `json:"name" binding:"max=50"`    // 诗人姓名
	DynastyID uint64 `json:"dynastyId"`                // 所属朝代ID
	Biography string `json:"biography"`                // 生平简介
	Avatar    string `json:"avatar" binding:"max=255"` // 头像URL
	BirthYear int    `json:"birthYear"`                // 出生年份
	DeathYear int    `json:"deathYear"`                // 逝世年份
}
