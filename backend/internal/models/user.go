package models

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
	Role     string `json:"role" gorm:"not null;default:'user';size:20;comment:角色"` // admin, user, reviewer, editor
	IsActive bool   `json:"isActive" gorm:"default:true;comment:是否启用"`
}

// TableName 指定表名为单数形式
func (User) TableName() string {
	return "user"
}

type Dynasty struct {
	BaseModel
	Name        string `json:"name" gorm:"not null;size:50;comment:朝代名称"`
	Period      string `json:"period" gorm:"size:100;comment:时期"`
	Description string `json:"description" gorm:"type:text;comment:描述"`
}

// TableName 指定表名为单数形式
func (Dynasty) TableName() string {
	return "dynasty"
}

type Poet struct {
	BaseModel
	Name      string  `json:"name" gorm:"not null;size:50;comment:诗人姓名"`
	DynastyID uint64  `json:"dynastyId" gorm:"comment:朝代ID"`
	Dynasty   Dynasty `json:"dynasty,omitempty"`
	Biography string  `json:"biography" gorm:"type:text;comment:生平简介"`
	Avatar    string  `json:"avatar" gorm:"size:255;comment:头像URL"`
	BirthYear int     `json:"birthYear" gorm:"comment:出生年份"`
	DeathYear int     `json:"deathYear" gorm:"comment:逝世年份"`
}

// TableName 指定表名为单数形式
func (Poet) TableName() string {
	return "poet"
}

type Poem struct {
	BaseModel
	Title        string  `json:"title" gorm:"not null;size:200;comment:诗词标题"`
	Content      string  `json:"content" gorm:"not null;type:text;comment:诗词内容"`
	AuthorID     uint64  `json:"authorId" gorm:"comment:作者ID"`
	Author       Poet    `json:"author,omitempty"`
	DynastyID    uint64  `json:"dynastyId" gorm:"comment:朝代ID"`
	Dynasty      Dynasty `json:"dynasty,omitempty"`
	Genre        string  `json:"genre" gorm:"size:50;comment:体裁"`
	Translation  string  `json:"translation" gorm:"type:text;comment:译文"`
	Appreciation string  `json:"appreciation" gorm:"type:text;comment:赏析"`
	Annotation   string  `json:"annotation" gorm:"type:text;comment:注释"`
	AudioURL     string  `json:"audioUrl" gorm:"size:500;comment:音频URL"`
	CoverImage   string  `json:"coverImage" gorm:"size:500;comment:封面图URL"`
	Views        int     `json:"views" gorm:"default:0;comment:浏览量"`
	Likes        int     `json:"likes" gorm:"default:0;comment:点赞数"`
	Dislikes     int     `json:"dislikes" gorm:"default:0;comment:点踩数"`
	Favorites    int     `json:"favorites" gorm:"default:0;comment:收藏数"`
}

// TableName 指定表名为单数形式
func (Poem) TableName() string {
	return "poem"
}
