package models

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
