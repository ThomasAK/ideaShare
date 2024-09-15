package models

type UserSetting struct {
	SoftDeleteModel
	Key    string `gorm:"not null" json:"key"`
	Value  string `gorm:"not null" json:"value"`
	UserId int    `gorm:"not null;index" json:"userId" faker:"-"`
}
