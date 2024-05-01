package model

type RoleUser struct {
	Id         *int    `json:"user_id" validate:"required" example:"1"`
	Email      *string `json:"email" validate:"required" example:"example@mail.ru"`
	UserRoleId *int    `json:"-"`
	RoleId     *int    `json:"-"`
}

type Role struct {
	Id    int        `json:"role_id" db:"role_id" example:"1" validate:"required"`
	Title string     `json:"title" db:"title" example:"ADMIN" validate:"required"`
	Users []RoleUser `json:"users" validate:"required"`
}

type UserRole struct {
	Id         *int    `json:"id" example:"1" validate:"required"`
	Title      *string `json:"title" example:"User" validate:"required"`
	UserId     *string `json:"-"`
	UserRoleId *int    `json:"-"`
}
