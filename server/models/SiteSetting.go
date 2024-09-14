package models

type SiteSetting struct {
	Base
	Name  string `json:"name"`
	Value string `json:"value"`
}
