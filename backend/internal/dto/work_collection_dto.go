package dto

import "time"

type WorkCollectionListRequest struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"pageSize"`
	Keyword   string `form:"keyword"`
	Published *bool  `form:"published"`
}

type WorkCollectionCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	CoverImage  string `json:"coverImage"`
	IsPublished bool   `json:"isPublished"`
}

type WorkCollectionUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverImage  string `json:"coverImage"`
	IsPublished bool   `json:"isPublished"`
}

type WorkCollectionItemCreateRequest struct {
	WorkType  string `json:"workType" binding:"required"`
	WorkID    uint64 `json:"workId" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

type WorkCollectionItemUpdateRequest struct {
	SortOrder int `json:"sortOrder"`
}

type WorkCollectionResponse struct {
	ID          uint64                       `json:"id"`
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	CoverImage  string                       `json:"coverImage"`
	ItemCount   int                          `json:"itemCount"`
	IsPublished bool                         `json:"isPublished"`
	Items       []WorkCollectionItemResponse `json:"items,omitempty"`
	CreatedAt   time.Time                    `json:"createdAt"`
	UpdatedAt   time.Time                    `json:"updatedAt"`
}

type WorkCollectionItemResponse struct {
	ID           uint64    `json:"id"`
	CollectionID uint64    `json:"collectionId"`
	WorkType     string    `json:"workType"`
	WorkID       uint64    `json:"workId"`
	SortOrder    int       `json:"sortOrder"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
