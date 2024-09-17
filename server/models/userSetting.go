package models

type UserSetting struct {
	SoftDeleteModel
	Key    string `gorm:"not null;index:user_settings_user_id_key_unq,unique;type:varchar(255)" json:"key"`
	Value  string `gorm:"not null" json:"value"`
	UserId int    `gorm:"not null;index:user_settings_user_id_key_unq,unique" json:"userId" faker:"-"`
}
