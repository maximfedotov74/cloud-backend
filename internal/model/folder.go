package model

import "time"

type Folder struct {
	FolderId       string     `json:"folder_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
	CreatedAt      time.Time  `json:"created_at" validate:"required"`
	UpdatedAt      time.Time  `json:"updated_at" validate:"required"`
	DeletedAt      *time.Time `json:"deleted_at"`
	Title          string     `json:"title" validate:"required" example:"folder1"`
	ParentFolderId *string    `json:"parent_folder_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6f"`
	UserId         string     `json:"user_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
}
