package models

type SiteSetting struct {
	SoftDeleteModel
	Key   string `gorm:"uniqueIndex;not null;type:varchar(255)" json:"name"`
	Value string `gorm:"not null" json:"value"`
}
