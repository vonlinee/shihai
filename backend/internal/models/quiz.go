package models

type Quiz struct {
	BaseModel
	PoemID        *uint64 `json:"poemId" gorm:"comment:诗词ID"`
	Poem          *Poem   `json:"poem,omitempty"`
	Question      string  `json:"question" gorm:"not null;type:text;comment:问题"`
	Options       string  `json:"options" gorm:"not null;type:json;comment:选项JSON"` // JSON array
	CorrectAnswer int     `json:"correctAnswer" gorm:"not null;comment:正确答案索引"`
	Explanation   string  `json:"explanation" gorm:"type:text;comment:答案解析"`
	Difficulty    string  `json:"difficulty" gorm:"default:'medium';size:10;comment:难度"` // easy, medium, hard
}

// TableName 指定表名为单数形式
func (Quiz) TableName() string {
	return "quiz"
}

type QuizRecord struct {
	BaseModel
	UserID    uint64 `json:"userId" gorm:"not null;index;comment:用户ID"`
	User      User   `json:"user,omitempty"`
	QuizID    uint64 `json:"quizId" gorm:"not null;index;comment:题目ID"`
	Quiz      Quiz   `json:"quiz,omitempty"`
	Answer    int    `json:"answer" gorm:"not null;comment:用户答案"`
	IsCorrect bool   `json:"isCorrect" gorm:"not null;comment:是否正确"`
}

// TableName 指定表名为单数形式
func (QuizRecord) TableName() string {
	return "quiz_record"
}
