package dto

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type PoemListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`
	Dynasty  string `form:"dynasty"`
	Author   string `form:"author"`
	Genre    string `form:"genre"`
}

type PoemCreateRequest struct {
	Title        string                         `json:"title" binding:"required,max=200"`
	Content      []string                       `json:"content" binding:"required"`
	AuthorID     RequestID                      `json:"authorId"`
	AuthorName   string                         `json:"authorName"`
	DynastyID    RequestID                      `json:"dynastyId"`
	DynastyName  string                         `json:"dynastyName"`
	Genre        string                         `json:"genre" binding:"max=50"`
	Translation  string                         `json:"translation"`
	Appreciation string                         `json:"appreciation"`
	Annotation   string                         `json:"annotation"`
	Annotations  *[]PoemAnnotationUpsertRequest `json:"annotations"`
	AudioURL     string                         `json:"audioUrl" binding:"max=500"`
	CoverImage   string                         `json:"coverImage" binding:"max=500"`
}

type PoemUpdateRequest struct {
	Title        string                         `json:"title" binding:"max=200"`
	Content      []string                       `json:"content"`
	AuthorID     RequestID                      `json:"authorId"`
	DynastyID    RequestID                      `json:"dynastyId"`
	Genre        string                         `json:"genre" binding:"max=50"`
	Translation  string                         `json:"translation"`
	Appreciation string                         `json:"appreciation"`
	Annotation   string                         `json:"annotation"`
	Annotations  *[]PoemAnnotationUpsertRequest `json:"annotations"`
	AudioURL     string                         `json:"audioUrl" binding:"max=500"`
	CoverImage   string                         `json:"coverImage" binding:"max=500"`
}

type BatchDeleteRequest struct {
	IDs IDList `json:"ids" binding:"required,min=1,dive,gt=0"`
}

type IDList []uint64

type RequestID uint64

func (ids *IDList) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	values := make([]uint64, 0, len(raw))
	for _, item := range raw {
		var stringID string
		if err := json.Unmarshal(item, &stringID); err == nil {
			id, parseErr := strconv.ParseUint(stringID, 10, 64)
			if parseErr != nil {
				return fmt.Errorf("invalid id %q", stringID)
			}
			values = append(values, id)
			continue
		}

		var numericID uint64
		if err := json.Unmarshal(item, &numericID); err != nil {
			return fmt.Errorf("invalid id")
		}
		values = append(values, numericID)
	}

	*ids = values
	return nil
}

func (id *RequestID) UnmarshalJSON(data []byte) error {
	value, err := parseJSONID(data)
	if err != nil {
		return err
	}
	*id = RequestID(value)
	return nil
}

func RequestIDPtrValue(id *RequestID) *uint64 {
	if id == nil {
		return nil
	}
	value := uint64(*id)
	return &value
}

func parseJSONID(data []byte) (uint64, error) {
	var stringID string
	if err := json.Unmarshal(data, &stringID); err == nil {
		id, parseErr := strconv.ParseUint(stringID, 10, 64)
		if parseErr != nil {
			return 0, fmt.Errorf("invalid id %q", stringID)
		}
		return id, nil
	}

	var numericID uint64
	if err := json.Unmarshal(data, &numericID); err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return numericID, nil
}

type PoemResponse struct {
	ID           uint64                   `json:"id,string"`
	Title        string                   `json:"title"`
	Content      []string                 `json:"content"`
	AuthorID     uint64                   `json:"authorId,string"`
	Author       AuthorResponse           `json:"author,omitempty"`
	DynastyID    uint64                   `json:"dynastyId,string"`
	Dynasty      DynastyResponse          `json:"dynasty,omitempty"`
	Genre        string                   `json:"genre"`
	Translation  string                   `json:"translation"`
	Appreciation string                   `json:"appreciation"`
	Annotation   string                   `json:"annotation"`
	Annotations  []PoemAnnotationResponse `json:"annotations,omitempty"`
	AudioURL     string                   `json:"audioUrl"`
	CoverImage   string                   `json:"coverImage"`
	Views        int                      `json:"views"`
	Likes        int                      `json:"likes"`
	Dislikes     int                      `json:"dislikes"`
	Favorites    int                      `json:"favorites"`
	CreatedAt    time.Time                `json:"createdAt"`
	UpdatedAt    time.Time                `json:"updatedAt"`
}

type AuthorResponse struct {
	ID        uint64 `json:"id,string"`
	Name      string `json:"name"`
	Biography string `json:"biography"`
	Avatar    string `json:"avatar"`
}

type PoetResponse struct {
	ID        uint64          `json:"id,string"`
	AuthorID  uint64          `json:"authorId,string"`
	Name      string          `json:"name"`
	DynastyID uint64          `json:"dynastyId,string"`
	Dynasty   DynastyResponse `json:"dynasty,omitempty"`
	Biography string          `json:"biography"`
	Avatar    string          `json:"avatar"`
	BirthYear int             `json:"birthYear"`
	DeathYear int             `json:"deathYear"`
}

type DynastyResponse struct {
	ID          uint64 `json:"id,string"`
	Name        string `json:"name"`
	Period      string `json:"period"`
	Description string `json:"description"`
}

type DynastyCreateRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Period      string `json:"period" binding:"max=100"`
	Description string `json:"description"`
}

type DynastyUpdateRequest struct {
	Name        string `json:"name" binding:"max=50"`
	Period      string `json:"period" binding:"max=100"`
	Description string `json:"description"`
}

type PoetCreateRequest struct {
	Name      string    `json:"name" binding:"required,max=50"`
	DynastyID RequestID `json:"dynastyId"`
	Biography string    `json:"biography"`
	Avatar    string    `json:"avatar" binding:"max=255"`
	BirthYear int       `json:"birthYear"`
	DeathYear int       `json:"deathYear"`
}

type PoetUpdateRequest struct {
	Name      string    `json:"name" binding:"max=50"`
	DynastyID RequestID `json:"dynastyId"`
	Biography string    `json:"biography"`
	Avatar    string    `json:"avatar" binding:"max=255"`
	BirthYear int       `json:"birthYear"`
	DeathYear int       `json:"deathYear"`
}
