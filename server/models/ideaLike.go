package models

import "gorm.io/gorm"

type IdeaLike struct {
	gorm.Model
	IdeaID    int
	CreatedBy int
}
