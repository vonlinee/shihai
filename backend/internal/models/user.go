package models

// User 用户表，存储平台所有注册用户的基本信息
// 用户通过 UserRole 关联表与角色绑定，实现 RBAC 权限控制
type User struct {
	BaseModel
	Username string `json:"username" gorm:"uniqueIndex;not null;size:50;comment:用户名"` // 用户名，唯一索引，用于登录
	Password string `json:"-" gorm:"not null;size:255;comment:密码"`                    // 密码，bcrypt 加密存储，JSON序列化时忽略
	Name     string `json:"name" gorm:"not null;size:50;comment:姓名"`                  // 用户显示名称
	Avatar   string `json:"avatar" gorm:"size:255;comment:头像URL"`                     // 头像图片地址
	Gender   string `json:"gender" gorm:"size:10;comment:性别"`                         // 性别：male/female/other
	Age      int    `json:"age" gorm:"comment:年龄"`                                    // 年龄
	Phone    string `json:"phone" gorm:"size:20;comment:手机号"`                         // 手机号码
	IDCard   string `json:"idCard" gorm:"size:18;comment:身份证号"`                       // 身份证号码
	IsActive bool   `json:"isActive" gorm:"default:true;comment:是否启用"`                // 是否启用，false 表示账户被禁用
}

// TableName 指定表名为单数形式
func (User) TableName() string {
	return "user"
}

// Dynasty 朝代表，存储中国历史朝代信息
// 诗人(Poet)通过 DynastyID 关联到具体朝代
type Dynasty struct {
	BaseModel
	Name        string `json:"name" gorm:"not null;size:50;comment:朝代名称"` // 朝代名称，如"唐"、"宋"
	Period      string `json:"period" gorm:"size:100;comment:时期"`         // 时期描述，如"618-907"
	Description string `json:"description" gorm:"type:text;comment:描述"`   // 朝代详细描述
}

// TableName 指定表名为单数形式
func (Dynasty) TableName() string {
	return "dynasty"
}

// Poet 诗人表，存储古诗词作者的基本信息
// 诗人通过 DynastyID 关联所属朝代，诗词(Poem)通过 AuthorID 关联诗人
type Poet struct {
	BaseModel
	Name      string  `json:"name" gorm:"not null;size:50;comment:诗人姓名"` // 诗人姓名
	DynastyID uint64  `json:"dynastyId" gorm:"comment:朝代ID"`             // 所属朝代ID，关联 Dynasty 表
	Dynasty   Dynasty `json:"dynasty,omitempty"`                         // 关联的朝代信息，查询时预加载
	Biography string  `json:"biography" gorm:"type:text;comment:生平简介"`   // 诗人生平简介
	Avatar    string  `json:"avatar" gorm:"size:255;comment:头像URL"`      // 诗人头像图片地址
	BirthYear int     `json:"birthYear" gorm:"comment:出生年份"`             // 出生年份，如701
	DeathYear int     `json:"deathYear" gorm:"comment:逝世年份"`             // 逝世年份，如762
}

// TableName 指定表名为单数形式
func (Poet) TableName() string {
	return "poet"
}

// Poem 诗词表，存储古诗词的核心内容及其元数据
// 诗词通过 AuthorID 关联诗人(Poet)，通过 DynastyID 关联朝代(Dynasty)
// 支持多种体裁(genre)、译文、赏析和注释等扩展内容
type Poem struct {
	BaseModel
	Title        string  `json:"title" gorm:"not null;size:200;comment:诗词标题"`    // 诗词标题，如"静夜思"
	Content      string  `json:"content" gorm:"not null;type:text;comment:诗词内容"` // 诗词正文内容
	AuthorID     uint64  `json:"authorId" gorm:"comment:作者ID"`                   // 作者ID，关联 Poet 表
	Author       Poet    `json:"author,omitempty"`                               // 关联的诗人信息，查询时预加载
	DynastyID    uint64  `json:"dynastyId" gorm:"comment:朝代ID"`                  // 朝代ID，关联 Dynasty 表
	Dynasty      Dynasty `json:"dynasty,omitempty"`                              // 关联的朝代信息，查询时预加载
	Genre        string  `json:"genre" gorm:"size:50;comment:体裁"`                // 诗词体裁，如"五言绝句"、"七言律诗"
	Translation  string  `json:"translation" gorm:"type:text;comment:译文"`        // 现代汉语译文
	Appreciation string  `json:"appreciation" gorm:"type:text;comment:赏析"`       // 诗词赏析与解读
	Annotation   string  `json:"annotation" gorm:"type:text;comment:注释"`         // 字词注释说明
	AudioURL     string  `json:"audioUrl" gorm:"size:500;comment:音频URL"`         // 诗词朗读音频地址
	CoverImage   string  `json:"coverImage" gorm:"size:500;comment:封面图URL"`      // 诗词封面图片地址
	Views        int     `json:"views" gorm:"default:0;comment:浏览量"`             // 浏览次数
	Likes        int     `json:"likes" gorm:"default:0;comment:点赞数"`             // 点赞数量
	Dislikes     int     `json:"dislikes" gorm:"default:0;comment:点踩数"`          // 点踩数量
	Favorites    int     `json:"favorites" gorm:"default:0;comment:收藏数"`         // 收藏数量
}

// TableName 指定表名为单数形式
func (Poem) TableName() string {
	return "poem"
}
