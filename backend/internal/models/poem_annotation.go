package models

const (
	PoemAnnotationTargetContent = "content"
	PoemAnnotationTypeNote      = "note"
)

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
