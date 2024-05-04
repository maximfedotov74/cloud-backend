package dto

type AddRoleToUser struct {
	Title  string `json:"title" validate:"required,min=3" example:"ADMIN"`
	UserId string `json:"user_id" validate:"required,uuid4" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d"`
}

type CreateRole struct {
	Title string `json:"title" validate:"required,min=6,max=55" db:"title" example:"ADMIN"`
}
