package models

type IdeaComment struct {
	SoftDeleteModel
	IdeaID  int    `gorm:"index;not null" json:"ideaID" faker:"-"`
	Comment string `gorm:"check:length(comment) > 3;not null" json:"comment"`
}
