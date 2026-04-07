package models

import (
	"shihai/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含通用审计字段
type BaseModel struct {
	ID        uint64         `json:"id" gorm:"primaryKey;comment:主键ID"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	CreatedBy uint64         `json:"createdBy" gorm:"default:0;comment:创建人ID"`
	UpdatedBy uint64         `json:"updatedBy" gorm:"default:0;comment:更新人ID"`
}

// BeforeCreate 创建前回调，自动生成雪花ID
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == 0 {
		m.ID = utils.GenerateID()
	}
	return nil
}
