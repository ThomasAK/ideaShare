package models

type User struct {
	Base
	ExternalID string     `json:"externalID"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	Roles      []UserRole `json:"roles"`
}
