package models

type Idea struct {
	Base
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	CreatedBy   int           `json:"created_by"`
	Likes       int           `gorm:"-" json:"likes"`
	Comments    []IdeaComment `json:"comments"`
	LikedByUser bool          `gorm:"-" json:"likedByUser"`
}
