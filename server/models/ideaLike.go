package models

type IdeaLike struct {
	Base
	IdeaID    int `json:"ideaID"`
	CreatedBy int `json:"createdBy"`
}
