package models

// Author stores common author identity shared by poems, poets, novels, and articles.
type Author struct {
	BaseModel
	Name      string `json:"name" gorm:"not null;size:50;uniqueIndex;comment:作者姓名"`
	Biography string `json:"biography" gorm:"type:text;comment:作者简介"`
	Avatar    string `json:"avatar" gorm:"size:255;comment:头像URL"`
}

// TableName specifies the database table name.
func (Author) TableName() string {
	return "author"
}
