package models

// WorkCollection 作品集表，存储一组可排序的作品条目
// 作品条目通过 WorkCollectionItem 使用多态关联，支持诗词、小说、文章等不同作品类型
type WorkCollection struct {
	BaseModel
	Title       string               `json:"title" gorm:"not null;size:200;comment:作品集标题"`   // 作品集标题
	Description string               `json:"description" gorm:"type:text;comment:作品集描述"`     // 作品集描述
	CoverImage  string               `json:"coverImage" gorm:"size:500;comment:封面图片URL"`     // 封面图片地址
	ItemCount   int                  `json:"itemCount" gorm:"default:0;comment:作品数量"`        // 作品条目数量
	IsPublished bool                 `json:"isPublished" gorm:"default:false;comment:是否发布"`  // 是否发布
	Items       []WorkCollectionItem `json:"items,omitempty" gorm:"foreignKey:CollectionID"` // 作品集条目
}

// TableName 指定表名为单数形式
func (WorkCollection) TableName() string {
	return "work_collection"
}

// WorkCollectionItem 作品集条目表，存储作品集与不同类型作品的关联关系
// WorkType 标识目标作品类型，如 poem、novel、article；WorkID 为对应作品表主键
type WorkCollectionItem struct {
	BaseModel
	CollectionID uint64         `json:"collectionId" gorm:"not null;index;uniqueIndex:idx_collection_work;comment:作品集ID"` // 作品集ID
	Collection   WorkCollection `json:"collection,omitempty"`                                                             // 所属作品集
	WorkType     string         `json:"workType" gorm:"not null;size:50;index:idx_work_ref;uniqueIndex:idx_collection_work;comment:作品类型"`
	WorkID       uint64         `json:"workId" gorm:"not null;index:idx_work_ref;uniqueIndex:idx_collection_work;comment:作品ID"`
	SortOrder    int            `json:"sortOrder" gorm:"default:0;comment:排序值"` // 作品集内排序值
}

// TableName 指定表名为单数形式
func (WorkCollectionItem) TableName() string {
	return "work_collection_item"
}
