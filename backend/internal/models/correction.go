package models

// CorrectionRequest 纠错申请表，存储用户提交的诗词内容纠错请求
// 纠错流程：提交(pending) → 投票中(voting) → 已批准(approved)/已拒绝(rejected) → 已完成(completed)
// 纠错类型覆盖诗词的各部分内容：正文、译文、赏析、注释
type CorrectionRequest struct {
	BaseModel
	PoemID        uint64 `json:"poemId" gorm:"not null;index;comment:诗词ID"`                   // 被纠错的诗词ID，关联 Poem 表
	Poem          Poem   `json:"poem,omitempty"`                                              // 关联的诗词信息，查询时预加载
	UserID        uint64 `json:"userId" gorm:"not null;index;comment:申请人ID"`                  // 纠错提交者ID，关联 User 表
	User          User   `json:"user,omitempty"`                                              // 关联的提交者信息，查询时预加载
	Type          string `json:"type" gorm:"not null;size:20;comment:纠错类型"`                   // 纠错类型：content(正文)/translation(译文)/appreciation(赏析)/annotation(注释)
	OriginalText  string `json:"originalText" gorm:"not null;type:text;comment:原文内容"`         // 被纠错的原文内容
	SuggestedText string `json:"suggestedText" gorm:"not null;type:text;comment:建议修改内容"`      // 建议修改后的内容
	Reason        string `json:"reason" gorm:"not null;type:text;comment:纠错理由"`               // 提交纠错的原因说明
	Status        string `json:"status" gorm:"not null;default:'pending';size:20;comment:状态"` // 当前状态：pending/voting/approved/rejected/completed
	VoteCount     int    `json:"voteCount" gorm:"default:0;comment:投票总数"`                     // 总投票数 = ApproveCount + RejectCount
	ApproveCount  int    `json:"approveCount" gorm:"default:0;comment:支持票数"`                  // 赞成票数量
	RejectCount   int    `json:"rejectCount" gorm:"default:0;comment:反对票数"`                   // 反对票数量
}

// TableName 指定表名为单数形式
func (CorrectionRequest) TableName() string {
	return "correction_request"
}

// CorrectionVote 纠错投票表，记录用户对纠错申请的投票
// 每个用户对同一纠错申请只能投一次票（approve 或 reject）
type CorrectionVote struct {
	BaseModel
	CorrectionID uint64            `json:"correctionId" gorm:"not null;index;comment:纠错申请ID"` // 被投票的纠错申请ID，关联 CorrectionRequest 表
	Correction   CorrectionRequest `json:"correction,omitempty"`                              // 关联的纠错申请信息，查询时预加载
	UserID       uint64            `json:"userId" gorm:"not null;index;comment:投票人ID"`        // 投票人ID，关联 User 表
	User         User              `json:"user,omitempty"`                                    // 关联的投票人信息，查询时预加载
	Type         string            `json:"type" gorm:"not null;size:10;comment:投票类型"`         // 投票类型：approve(赞成) 或 reject(反对)
	Comment      string            `json:"comment" gorm:"type:text;comment:投票意见"`             // 投票附带的文字意见
}

// TableName 指定表名为单数形式
func (CorrectionVote) TableName() string {
	return "correction_vote"
}
