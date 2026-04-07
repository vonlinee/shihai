package models

type CorrectionRequest struct {
	BaseModel
	PoemID        uint64 `json:"poemId" gorm:"not null;index;comment:诗词ID"`
	Poem          Poem   `json:"poem,omitempty"`
	UserID        uint64 `json:"userId" gorm:"not null;index;comment:申请人ID"`
	User          User   `json:"user,omitempty"`
	Type          string `json:"type" gorm:"not null;size:20;comment:纠错类型"` // content, translation, appreciation, annotation
	OriginalText  string `json:"originalText" gorm:"not null;type:text;comment:原文内容"`
	SuggestedText string `json:"suggestedText" gorm:"not null;type:text;comment:建议修改内容"`
	Reason        string `json:"reason" gorm:"not null;type:text;comment:纠错理由"`
	Status        string `json:"status" gorm:"not null;default:'pending';size:20;comment:状态"` // pending, voting, approved, rejected, completed
	VoteCount     int    `json:"voteCount" gorm:"default:0;comment:投票总数"`
	ApproveCount  int    `json:"approveCount" gorm:"default:0;comment:支持票数"`
	RejectCount   int    `json:"rejectCount" gorm:"default:0;comment:反对票数"`
}

// TableName 指定表名为单数形式
func (CorrectionRequest) TableName() string {
	return "correction_request"
}

type CorrectionVote struct {
	BaseModel
	CorrectionID uint64          `json:"correctionId" gorm:"not null;index;comment:纠错申请ID"`
	Correction   CorrectionRequest `json:"correction,omitempty"`
	UserID       uint64          `json:"userId" gorm:"not null;index;comment:投票人ID"`
	User         User            `json:"user,omitempty"`
	Type         string          `json:"type" gorm:"not null;size:10;comment:投票类型"` // approve, reject
	Comment      string          `json:"comment" gorm:"type:text;comment:投票意见"`
}

// TableName 指定表名为单数形式
func (CorrectionVote) TableName() string {
	return "correction_vote"
}
