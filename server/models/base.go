package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel interface {
	comparable
	GetID() int
	SetCreatedBy(userId int)
	GetCreatedBy() int
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetUpdatedAt(time time.Time)
}

type HardDeleteModel struct {
	ID        int       `gorm:"primarykey" json:"id,omitempty" faker:"-"`
	CreatedAt time.Time `gorm:"index;not null" json:"createdAt" faker:"-"`
	UpdatedAt time.Time `gorm:"index" json:"updatedAt" faker:"-"`
	CreatedBy int       `gorm:"index;not null" json:"createdBy" faker:"-"`
}

type SoftDeleteModel struct {
	HardDeleteModel
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" faker:"-"`
}

func (b *HardDeleteModel) GetID() int {
	return b.ID
}

func (b *HardDeleteModel) SetCreatedBy(userId int) {
	b.CreatedBy = userId
}

func (b *HardDeleteModel) GetCreatedBy() int {
	return b.CreatedBy
}

func (b *HardDeleteModel) GetCreatedAt() time.Time {
	return b.CreatedAt
}

func (b *HardDeleteModel) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}

func (b *HardDeleteModel) SetUpdatedAt(time time.Time) {
	b.UpdatedAt = time
}
