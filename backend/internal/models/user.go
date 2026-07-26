package models

// User stores platform account information.
type User struct {
	BaseModel
	Username string `json:"username" gorm:"uniqueIndex;not null;size:50;comment:用户名"`
	Password string `json:"-" gorm:"not null;size:255;comment:密码"`
	Name     string `json:"name" gorm:"not null;size:50;comment:姓名"`
	Avatar   string `json:"avatar" gorm:"size:255;comment:头像URL"`
	Gender   string `json:"gender" gorm:"size:10;comment:性别"`
	Age      int    `json:"age" gorm:"comment:年龄"`
	Phone    string `json:"phone" gorm:"size:20;comment:手机号"`
	IDCard   string `json:"idCard" gorm:"size:18;comment:身份证号"`
	IsActive bool   `json:"isActive" gorm:"default:true;comment:是否启用"`
}

// TableName specifies the database table name.
func (User) TableName() string {
	return "user"
}
