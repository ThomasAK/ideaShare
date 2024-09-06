package models

import "gorm.io/gorm"

type IdeaComment struct {
	gorm.Model
	IdeaID    int
	Comment   string
	CreatedBy int
}
