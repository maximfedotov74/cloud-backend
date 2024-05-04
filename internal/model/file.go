package model

import "time"

type File struct {
	FileId    string     `json:"file_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
	CreatedAt time.Time  `json:"created_at" validate:"required"`
	UpdatedAt time.Time  `json:"updated_at" validate:"required"`
	DeletedAt *time.Time `json:"deleted_at"`
	Title     string     `json:"title" validate:"required" example:"file1"`
	Ext       string     `json:"ext" validate:"required" example:".png"`
	Size      int64      `json:"size" validate:"required" example:"123"`
	FolderId  string     `json:"folder_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
	UserId    string     `json:"user_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
}
