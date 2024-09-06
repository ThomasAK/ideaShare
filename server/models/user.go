package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ExternalID string
	FirstName  string
	LastName   string
	Roles      []UserRole
}
