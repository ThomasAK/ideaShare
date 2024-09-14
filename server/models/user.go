package models

type User struct {
	Base
	ExternalID string         `gorm:"index" json:"externalID"`
	FirstName  string         `json:"firstName"`
	LastName   string         `json:"lastName"`
	Roles      []*UserRole    `json:"roles"`
	Settings   []*UserSetting `json:"settings"`
	Ideas      []*Idea        `gorm:"foreignKey:CreatedBy" json:"ideas"`
	Likes      []*IdeaLike    `gorm:"foreignKey:CreatedBy" json:"likes"`
	Comments   []*IdeaComment `gorm:"foreignKey:CreatedBy" json:"comments"`
}
