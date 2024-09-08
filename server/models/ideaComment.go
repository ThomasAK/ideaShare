package models

type IdeaComment struct {
	Base
	IdeaID    int    `json:"ideaID"`
	Comment   string `json:"comment"`
	CreatedBy int    `json:"createdBy"`
}
