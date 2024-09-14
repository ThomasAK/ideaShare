package models

type Role = string

const (
	SiteAdmin Role = "siteAdmin"
	IdeaAdmin Role = "ideaAdmin"
)

type UserRole struct {
	Base
	UserID int  `json:"userID"`
	Role   Role `json:"role"`
}
