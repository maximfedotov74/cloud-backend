package dto

type CreateFile struct {
	Title    string `json:"title" example:"folder1" validate:"required,min=1"`
	Ext      string `json:"ext" validate:"required" example:".png"`
	Size     int64  `json:"size" validate:"required" example:"123"`
	FolderId string `json:"folder_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
	UserId   string `json:"user_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
}
