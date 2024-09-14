package models

type IdeaComment struct {
	Base
	IdeaID  int    `gorm:"index" json:"ideaID"`
	Comment string `gorm:"check:length(comment) > 3" json:"comment"`
}
