package model

type Hello struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type CreateMsgDto struct {
	Message string `json:"message" validate:"required,min=5"`
}
