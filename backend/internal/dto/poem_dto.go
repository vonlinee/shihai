package dto

import "time"

// PoemListRequest 诗词列表请求
type PoemListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`
	Dynasty  string `form:"dynasty"`
	Author   string `form:"author"`
	Genre    string `form:"genre"`
}

// PoemCreateRequest 创建诗词请求
type PoemCreateRequest struct {
	Title        string `json:"title" binding:"required,max=200"`
	Content      string `json:"content" binding:"required"`
	AuthorID     uint64 `json:"authorId" binding:"required"`
	DynastyID    uint64 `json:"dynastyId" binding:"required"`
	Genre        string `json:"genre" binding:"max=50"`
	Translation  string `json:"translation"`
	Appreciation string `json:"appreciation"`
	Annotation   string `json:"annotation"`
	AudioURL     string `json:"audioUrl" binding:"max=500"`
	CoverImage   string `json:"coverImage" binding:"max=500"`
}

// PoemUpdateRequest 更新诗词请求
type PoemUpdateRequest struct {
	Title        string `json:"title" binding:"max=200"`
	Content      string `json:"content"`
	AuthorID     uint64 `json:"authorId"`
	DynastyID    uint64 `json:"dynastyId"`
	Genre        string `json:"genre" binding:"max=50"`
	Translation  string `json:"translation"`
	Appreciation string `json:"appreciation"`
	Annotation   string `json:"annotation"`
	AudioURL     string `json:"audioUrl" binding:"max=500"`
	CoverImage   string `json:"coverImage" binding:"max=500"`
}

// PoemResponse 诗词响应
type PoemResponse struct {
	ID           uint64          `json:"id"`
	Title        string          `json:"title"`
	Content      string          `json:"content"`
	AuthorID     uint64          `json:"authorId"`
	Author       PoetResponse    `json:"author,omitempty"`
	DynastyID    uint64          `json:"dynastyId"`
	Dynasty      DynastyResponse `json:"dynasty,omitempty"`
	Genre        string          `json:"genre"`
	Translation  string          `json:"translation"`
	Appreciation string          `json:"appreciation"`
	Annotation   string          `json:"annotation"`
	AudioURL     string          `json:"audioUrl"`
	CoverImage   string          `json:"coverImage"`
	Views        int             `json:"views"`
	Likes        int             `json:"likes"`
	Dislikes     int             `json:"dislikes"`
	Favorites    int             `json:"favorites"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

// PoetResponse 诗人响应
type PoetResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	DynastyID uint64 `json:"dynastyId"`
	Biography string `json:"biography"`
	Avatar    string `json:"avatar"`
}

// DynastyResponse 朝代响应
type DynastyResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Period      string `json:"period"`
	Description string `json:"description"`
}

// DynastyCreateRequest 创建朝代请求
type DynastyCreateRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Period      string `json:"period" binding:"max=100"`
	Description string `json:"description"`
}

// PoetCreateRequest 创建诗人请求
type PoetCreateRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	DynastyID uint64 `json:"dynastyId" binding:"required"`
	Biography string `json:"biography"`
	Avatar    string `json:"avatar" binding:"max=255"`
	BirthYear int    `json:"birthYear"`
	DeathYear int    `json:"deathYear"`
}
