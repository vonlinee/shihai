package models

const (
	PoemAnnotationTargetContent = "content"
	PoemAnnotationTypeNote      = "note"
)

// Dynasty stores historical dynasty information.
type Dynasty struct {
	BaseModel
	Name        string `json:"name" gorm:"not null;size:50;comment:朝代名称"`
	NameEn      string `json:"nameEn" gorm:"size:100;comment:朝代英文名称"`
	StartYear   *int   `json:"startYear" gorm:"comment:朝代开始年份"`
	EndYear     *int   `json:"endYear" gorm:"comment:朝代结束年份"`
	Period      string `json:"period" gorm:"size:100;comment:时期"`
	Description string `json:"description" gorm:"type:text;comment:描述"`
}

// TableName specifies the database table name.
func (Dynasty) TableName() string {
	return "dynasty"
}

// Poet stores poetry-domain metadata for an Author.
type Poet struct {
	BaseModel
	AuthorID  uint64  `json:"authorId" gorm:"index;comment:作者ID"`
	Author    Author  `json:"author,omitempty"`
	DynastyID uint64  `json:"dynastyId" gorm:"comment:朝代ID"`
	Dynasty   Dynasty `json:"dynasty,omitempty"`
	BirthYear int     `json:"birthYear" gorm:"comment:出生年份"`
	DeathYear int     `json:"deathYear" gorm:"comment:逝世年份"`
}

// TableName specifies the database table name.
func (Poet) TableName() string {
	return "poet"
}

// Poem stores poem content and metadata.
type Poem struct {
	BaseModel
	Title         string   `json:"title" gorm:"not null;size:500;comment:诗词标题"`
	Content       []string `json:"content" gorm:"not null;serializer:json;type:jsonb;comment:诗词内容"`
	Pingze        []string `json:"pingze" gorm:"serializer:json;type:jsonb;comment:诗词正文平仄，按正文文本行存储平/仄标记"` // Pingze 按 Content 每个文本元素对应存储平仄标记。
	AuthorID      uint64   `json:"authorId" gorm:"comment:作者ID"`
	Author        Author   `json:"author,omitempty"`
	DynastyID     uint64   `json:"dynastyId" gorm:"comment:朝代ID"`
	Dynasty       Dynasty  `json:"dynasty,omitempty"`
	GenreCategory string   `json:"genreCategory" gorm:"size:50;comment:体裁一级分类"`
	Genre         string   `json:"genre" gorm:"size:50;comment:体裁"`
	CiTuneID      *uint64  `json:"ciTuneId" gorm:"index;comment:词牌ID，仅词使用"`
	CiTune        *CiTune  `json:"ciTune,omitempty" gorm:"foreignKey:CiTuneID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Translation   string   `json:"translation" gorm:"type:text;comment:译文"`
	Appreciation  string   `json:"appreciation" gorm:"type:text;comment:赏析"`
	Annotation    string   `json:"annotation" gorm:"type:text;comment:注释"`
	AudioURL      string   `json:"audioUrl" gorm:"size:500;comment:音频URL"`
	CoverImage    string   `json:"coverImage" gorm:"size:500;comment:封面图URL"`
	Views         int      `json:"views" gorm:"default:0;comment:浏览量"`
	Likes         int      `json:"likes" gorm:"default:0;comment:点赞数"`
	Dislikes      int      `json:"dislikes" gorm:"default:0;comment:点踩数"`
	Favorites     int      `json:"favorites" gorm:"default:0;comment:收藏数"`
}

// TableName specifies the database table name.
func (Poem) TableName() string {
	return "poem"
}

// CiTune stores tune-pattern metadata for ci works.
type CiTune struct {
	BaseModel
	Name        string    `json:"name" gorm:"not null;size:100;uniqueIndex;comment:词牌名"`
	Aliases     []string  `json:"aliases" gorm:"serializer:json;type:jsonb;comment:词牌别名，用于同调异名和标题解析"`
	PoemTypeID  *uint64   `json:"poemTypeId" gorm:"index;comment:所属词体裁ID，例如小令、中调、长调"`
	PoemType    *PoemType `json:"poemType,omitempty" gorm:"foreignKey:PoemTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Description string    `json:"description" gorm:"type:text;comment:词牌说明"`
}

// TableName specifies the database table name.
func (CiTune) TableName() string {
	return "ci_tune"
}

// PoemAnnotation stores administrator-maintained notes for a selected poem text range.
type PoemAnnotation struct {
	BaseModel
	PoemID       uint64 `json:"poemId" gorm:"not null;index;comment:诗词ID"`
	Poem         Poem   `json:"poem,omitempty"`
	TargetField  string `json:"targetField" gorm:"not null;size:50;default:content;index;comment:标注目标字段"`
	StartLine    int    `json:"startLine" gorm:"not null;comment:起始行下标"`
	StartOffset  int    `json:"startOffset" gorm:"not null;comment:起始字符偏移"`
	EndLine      int    `json:"endLine" gorm:"not null;comment:结束行下标"`
	EndOffset    int    `json:"endOffset" gorm:"not null;comment:结束字符偏移"`
	SelectedText string `json:"selectedText" gorm:"not null;type:text;comment:被标注文案"`
	Title        string `json:"title" gorm:"size:100;comment:标注标题"`
	Content      string `json:"content" gorm:"not null;type:text;comment:标注内容"`
	Type         string `json:"type" gorm:"not null;size:50;default:note;comment:标注类型"`
	DisplayOrder int    `json:"displayOrder" gorm:"not null;default:0;comment:展示排序"`
}

// TableName specifies the database table name.
func (PoemAnnotation) TableName() string {
	return "poem_annotation"
}

// PoemType stores poem form metadata such as category, line count, and line length.
type PoemType struct {
	BaseModel
	Name         string `json:"name" gorm:"not null;size:50;uniqueIndex;comment:诗词类型名称"` // 诗词类型名称
	Category     string `json:"category" gorm:"not null;size:50;comment:诗词类型分类"`         // 诗词类型分类
	Lines        *int   `json:"lines" gorm:"comment:句数"`                                 // 句数，空值表示不限
	CharsPerLine *int   `json:"charsPerLine" gorm:"comment:每句字数"`                        // 每句字数，空值表示不限
	Description  string `json:"description" gorm:"type:text;comment:类型描述"`               // 类型描述
}

// TableName specifies the database table name.
func (PoemType) TableName() string {
	return "poem_type"
}
