package models

import "gorm.io/gorm"

type Idea struct {
	gorm.Model
	Title       string
	Description string
	Status      string
	CreatedBy   int
	Likes       int `gorm:"-"`
	Comments    []IdeaComment
}
