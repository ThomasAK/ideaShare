package models

type Role = string

const (
	SiteAdmin Role = "siteAdmin"
	IdeaAdmin Role = "ideaAdmin"
)

type UserRole struct {
	SoftDeleteModel
	UserID int  `gorm:"not null" json:"userID" faker:"-"`
	Role   Role `gorm:"not null" json:"role"`
}
