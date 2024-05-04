package dto

type CreateFolder struct {
	Title          string  `json:"title" example:"folder1" validate:"required,min=1"`
	ParentFolderId *string `json:"parent_folder_id" validate:"omitempty,uuid4" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d"`
}
