package models

type User struct {
	SoftDeleteModel
	ExternalID string         `gorm:"index:unique;not null" json:"externalID"`
	FirstName  string         `gorm:"not null" json:"firstName"`
	LastName   string         `json:"lastName"`
	Email      string         `gorm:"index:unique;not null" json:"email"`
	Roles      []*UserRole    `json:"roles" faker:"-"`
	Settings   []*UserSetting `json:"settings" faker:"-"`
}
