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
	Title         string                         `json:"title" binding:"required,max=200"`
	Content       []string                       `json:"content" binding:"required"`
	Pingze        []string                       `json:"pingze"` // Pingze 与 Content 每个文本元素对应的平仄标记，只允许平/仄/?。
	AuthorID      RequestID                      `json:"authorId"`
	AuthorName    string                         `json:"authorName"`
	DynastyID     RequestID                      `json:"dynastyId"`
	DynastyName   string                         `json:"dynastyName"`
	GenreCategory string                         `json:"genreCategory" binding:"max=50"` // GenreCategory 体裁一级分类，例如诗、词、曲、文。
	Genre         string                         `json:"genre" binding:"max=50"`
	CiTuneID      RequestID                      `json:"ciTuneId"` // CiTuneID 词牌 ID，仅当体裁一级分类为词时使用；未提供时可由标题解析。
	Translation   string                         `json:"translation"`
	Appreciation  string                         `json:"appreciation"`
	Annotation    string                         `json:"annotation"`
	Annotations   *[]PoemAnnotationUpsertRequest `json:"annotations"`
	AudioURL      string                         `json:"audioUrl" binding:"max=500"`
	CoverImage    string                         `json:"coverImage" binding:"max=500"`
}

type PoemUpdateRequest struct {
	Title         string                         `json:"title" binding:"max=200"`
	Content       []string                       `json:"content"`
	Pingze        *[]string                      `json:"pingze"` // Pingze 与 Content 每个文本元素对应的平仄标记，只允许平/仄/?；nil 表示本次不主动覆盖。
	AuthorID      RequestID                      `json:"authorId"`
	DynastyID     RequestID                      `json:"dynastyId"`
	GenreCategory string                         `json:"genreCategory" binding:"max=50"` // GenreCategory 体裁一级分类，例如诗、词、曲、文。
	Genre         string                         `json:"genre" binding:"max=50"`
	CiTuneID      *RequestID                     `json:"ciTuneId"` // CiTuneID 词牌 ID，nil 表示不主动覆盖，0 表示清空。
	Translation   string                         `json:"translation"`
	Appreciation  string                         `json:"appreciation"`
	Annotation    string                         `json:"annotation"`
	Annotations   *[]PoemAnnotationUpsertRequest `json:"annotations"`
	AudioURL      string                         `json:"audioUrl" binding:"max=500"`
	CoverImage    string                         `json:"coverImage" binding:"max=500"`
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
	ID            uint64                   `json:"id,string"`
	Title         string                   `json:"title"`
	Content       []string                 `json:"content"`
	Pingze        []string                 `json:"pingze"` // Pingze 与 Content 每个文本元素对应的平仄标记。
	AuthorID      uint64                   `json:"authorId,string"`
	Author        AuthorResponse           `json:"author,omitempty"`
	DynastyID     uint64                   `json:"dynastyId,string"`
	Dynasty       DynastyResponse          `json:"dynasty,omitempty"`
	GenreCategory string                   `json:"genreCategory"` // GenreCategory 体裁一级分类，例如诗、词、曲、文。
	Genre         string                   `json:"genre"`
	CiTuneID      *uint64                  `json:"ciTuneId,string,omitempty"` // CiTuneID 词牌 ID，仅词可能存在。
	CiTune        *CiTuneResponse          `json:"ciTune,omitempty"`          // CiTune 词牌摘要，仅词可能存在。
	Translation   string                   `json:"translation"`
	Appreciation  string                   `json:"appreciation"`
	Annotation    string                   `json:"annotation"`
	Annotations   []PoemAnnotationResponse `json:"annotations,omitempty"`
	AudioURL      string                   `json:"audioUrl"`
	CoverImage    string                   `json:"coverImage"`
	Views         int                      `json:"views"`
	Likes         int                      `json:"likes"`
	Dislikes      int                      `json:"dislikes"`
	Favorites     int                      `json:"favorites"`
	CreatedAt     time.Time                `json:"createdAt"`
	UpdatedAt     time.Time                `json:"updatedAt"`
}

// GenreCategoryResponse describes a top-level poem genre category and its selectable child genres.
type GenreCategoryResponse struct {
	Name   string          `json:"name"`   // Name 体裁大类名称，例如诗、词、曲、文。
	Genres []GenreResponse `json:"genres"` // Genres 当前大类下的细分类别列表。
}

// GenreResponse describes a selectable poem genre stored on Poem.Genre.
type GenreResponse struct {
	Name         string `json:"name"`         // Name 细分类别名称，也是诗词保存到 genre 字段的值。
	Lines        *int   `json:"lines"`        // Lines 常见句数；nil 表示不限定。
	CharsPerLine *int   `json:"charsPerLine"` // CharsPerLine 常见每句字数；nil 表示不限定。
	Description  string `json:"description"`  // Description 细分类别说明。
}

// PoemTypeResponse describes poem type reference data for admin management.
type PoemTypeResponse struct {
	ID           uint64    `json:"id,string"`    // ID 体裁的 Snowflake 主键，对外序列化为字符串。
	Name         string    `json:"name"`         // Name 细分类别名称，也是诗词保存到 genre 字段的值。
	Category     string    `json:"category"`     // Category 一级分类，例如诗、词、曲、文。
	Lines        *int      `json:"lines"`        // Lines 常见句数；nil 表示不限定。
	CharsPerLine *int      `json:"charsPerLine"` // CharsPerLine 常见每句字数；nil 表示不限定。
	Description  string    `json:"description"`  // Description 体裁说明。
	CreatedAt    time.Time `json:"createdAt"`    // CreatedAt 创建时间。
	UpdatedAt    time.Time `json:"updatedAt"`    // UpdatedAt 更新时间。
}

// PoemTypeCreateRequest describes the payload for creating poem type reference data.
type PoemTypeCreateRequest struct {
	Name         string `json:"name" binding:"required,max=50"`     // Name 细分类别名称，也是诗词保存到 genre 字段的值。
	Category     string `json:"category" binding:"required,max=50"` // Category 一级分类，例如诗、词、曲、文。
	Lines        *int   `json:"lines"`                              // Lines 常见句数；nil 表示不限定。
	CharsPerLine *int   `json:"charsPerLine"`                       // CharsPerLine 常见每句字数；nil 表示不限定。
	Description  string `json:"description"`                        // Description 体裁说明。
}

// PoemTypeUpdateRequest describes the full replacement payload for poem type reference data.
type PoemTypeUpdateRequest struct {
	Name         string `json:"name" binding:"required,max=50"`     // Name 细分类别名称，也是诗词保存到 genre 字段的值。
	Category     string `json:"category" binding:"required,max=50"` // Category 一级分类，例如诗、词、曲、文。
	Lines        *int   `json:"lines"`                              // Lines 常见句数；nil 表示不限定。
	CharsPerLine *int   `json:"charsPerLine"`                       // CharsPerLine 常见每句字数；nil 表示不限定。
	Description  string `json:"description"`                        // Description 体裁说明。
}

// CiTuneResponse describes ci tune reference data for admin management and poem editing.
type CiTuneResponse struct {
	ID          uint64            `json:"id,string"`                   // ID 词牌的 Snowflake 主键，对外序列化为字符串。
	Name        string            `json:"name"`                        // Name 词牌名，例如水调歌头、念奴娇。
	Aliases     []string          `json:"aliases"`                     // Aliases 词牌别名，用于兼容同调异名和标题解析。
	PoemTypeID  *uint64           `json:"poemTypeId,string,omitempty"` // PoemTypeID 所属词体裁 ID，例如小令、中调、长调；nil 表示未归类。
	PoemType    *PoemTypeResponse `json:"poemType,omitempty"`          // PoemType 所属词体裁摘要。
	Description string            `json:"description"`                 // Description 词牌说明。
	CreatedAt   time.Time         `json:"createdAt"`                   // CreatedAt 创建时间。
	UpdatedAt   time.Time         `json:"updatedAt"`                   // UpdatedAt 更新时间。
}

// CiTuneCreateRequest describes the payload for creating ci tune reference data.
type CiTuneCreateRequest struct {
	Name        string    `json:"name" binding:"required,max=100"` // Name 词牌名，例如水调歌头、念奴娇。
	Aliases     []string  `json:"aliases"`                         // Aliases 词牌别名，用于标题解析。
	PoemTypeID  RequestID `json:"poemTypeId"`                      // PoemTypeID 所属词体裁 ID；0 表示未归类。
	Description string    `json:"description"`                     // Description 词牌说明。
}

// CiTuneUpdateRequest describes the full replacement payload for ci tune reference data.
type CiTuneUpdateRequest struct {
	Name        string    `json:"name" binding:"required,max=100"` // Name 词牌名，例如水调歌头、念奴娇。
	Aliases     []string  `json:"aliases"`                         // Aliases 词牌别名，用于标题解析。
	PoemTypeID  RequestID `json:"poemTypeId"`                      // PoemTypeID 所属词体裁 ID；0 表示未归类。
	Description string    `json:"description"`                     // Description 词牌说明。
}

// CiTuneTitleParseRequest describes a title to parse into a ci tune.
type CiTuneTitleParseRequest struct {
	Title string `json:"title" binding:"required,max=200"` // Title 诗词标题，通常以词牌名开头。
}

// CiTuneTitleParseResponse describes the matched ci tune and matched text.
type CiTuneTitleParseResponse struct {
	CiTune      *CiTuneResponse `json:"ciTune,omitempty"` // CiTune 匹配到的词牌；未匹配时为 nil。
	MatchedName string          `json:"matchedName"`      // MatchedName 实际命中的词牌名或别名。
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
