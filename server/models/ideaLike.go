package models

type IdeaLike struct {
	HardDeleteModel
	IdeaID int `gorm:"not null;index:user_like_unq,unique" json:"ideaID" faker:"-"`
	UserID int `gorm:"not null;index:user_like_unq,unique" json:"userID" faker:"-"`
}
