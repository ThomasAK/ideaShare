package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel interface {
	GetID() int
	SetCreatedBy(userId int)
	GetCreatedBy() int
}

type Base struct {
	ID        int            `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"index" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	CreatedBy int            `gorm:"index" json:"createdBy"`
}

func (b *Base) GetID() int {
	return b.ID
}

func (b *Base) SetCreatedBy(userId int) {
	b.CreatedBy = userId
}

func (b *Base) GetCreatedBy() int {
	return b.CreatedBy
}
