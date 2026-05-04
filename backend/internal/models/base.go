package models

import (
	"shihai/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含通用审计字段
// 所有业务模型都应内嵌此结构体，以获得主键、时间戳和软删除等通用能力
// 主键 ID 采用雪花算法生成，确保分布式环境下的唯一性
type BaseModel struct {
	ID        uint64         `json:"id" gorm:"primaryKey;comment:主键ID"`        // 主键ID，雪花算法生成
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`            // 记录创建时间
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:更新时间"`            // 记录最后更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`              // 软删除时间，非空表示已删除，JSON序列化时忽略
	CreatedBy uint64         `json:"createdBy" gorm:"default:0;comment:创建人ID"` // 创建人ID，0表示系统创建
	UpdatedBy uint64         `json:"updatedBy" gorm:"default:0;comment:更新人ID"` // 最后更新人ID，0表示系统更新
}

// BeforeCreate 创建前回调，自动生成雪花ID
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == 0 {
		m.ID = utils.GenerateID()
	}
	return nil
}
