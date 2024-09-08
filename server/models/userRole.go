package models

type UserRole struct {
	Base
	UserID int    `json:"userID"`
	Role   string `json:"role"`
}
