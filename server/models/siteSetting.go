package models

type SiteSetting struct {
	SoftDeleteModel
	Key   string `gorm:"not null" json:"name"`
	Value string `gorm:"not null" json:"value"`
}
