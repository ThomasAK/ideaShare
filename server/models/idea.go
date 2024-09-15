package models

type Idea struct {
	SoftDeleteModel
	Title       string         `gorm:"index;not null" json:"title"`
	Description string         `json:"description"`
	Status      string         `gorm:"index;not null" json:"status"`
	Likes       int            `gorm:"-" json:"likes" faker:"-"`
	Comments    []*IdeaComment `json:"comments" faker:"-"`
	LikedByUser bool           `gorm:"-" json:"likedByUser" faker:"-"`
}
