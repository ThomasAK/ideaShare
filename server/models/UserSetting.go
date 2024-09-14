package models

type UserSetting struct {
	Base
	Key    string `json:"key"`
	Value  string `json:"value"`
	UserId int    `json:"userId"`
}
