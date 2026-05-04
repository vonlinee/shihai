package models

// Quiz 问答题目表，存储古诗词相关的问答题目
// 题目可关联具体诗词(PoemID)，也可为独立题目（PoemID 为 nil）
// 选项以 JSON 数组格式存储，CorrectAnswer 为正确选项的索引
type Quiz struct {
	BaseModel
	PoemID        *uint64 `json:"poemId" gorm:"comment:诗词ID"`                            // 关联诗词ID，可选，关联 Poem 表
	Poem          *Poem   `json:"poem,omitempty"`                                        // 关联的诗词信息，查询时预加载
	Question      string  `json:"question" gorm:"not null;type:text;comment:问题"`         // 题目内容
	Options       string  `json:"options" gorm:"not null;type:json;comment:选项JSON"`      // 选项列表，JSON 数组格式，如 ["A选项","B选项"]
	CorrectAnswer int     `json:"correctAnswer" gorm:"not null;comment:正确答案索引"`          // 正确选项的索引，从 0 开始
	Explanation   string  `json:"explanation" gorm:"type:text;comment:答案解析"`             // 答案解析说明
	Difficulty    string  `json:"difficulty" gorm:"default:'medium';size:10;comment:难度"` // 题目难度：easy(简单)/medium(中等)/hard(困难)
}

// TableName 指定表名为单数形式
func (Quiz) TableName() string {
	return "quiz"
}

// QuizRecord 答题记录表，存储用户的答题历史
// 用于统计用户答题正确率和学习进度
type QuizRecord struct {
	BaseModel
	UserID    uint64 `json:"userId" gorm:"not null;index;comment:用户ID"` // 答题用户ID，关联 User 表
	User      User   `json:"user,omitempty"`                            // 关联的用户信息，查询时预加载
	QuizID    uint64 `json:"quizId" gorm:"not null;index;comment:题目ID"` // 答题的题目ID，关联 Quiz 表
	Quiz      Quiz   `json:"quiz,omitempty"`                            // 关联的题目信息，查询时预加载
	Answer    int    `json:"answer" gorm:"not null;comment:用户答案"`       // 用户选择的答案索引，从 0 开始
	IsCorrect bool   `json:"isCorrect" gorm:"not null;comment:是否正确"`    // 是否回答正确
}

// TableName 指定表名为单数形式
func (QuizRecord) TableName() string {
	return "quiz_record"
}
