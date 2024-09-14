package models

type Idea struct {
	Base
	Title       string         `gorm:"index" json:"title"`
	Description string         `json:"description"`
	Status      string         `gorm:"index" json:"status"`
	Likes       int            `gorm:"-" json:"likes"`
	Comments    []*IdeaComment `json:"comments"`
	LikedByUser bool           `gorm:"-" json:"likedByUser"`
}
