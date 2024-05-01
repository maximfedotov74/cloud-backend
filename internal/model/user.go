package model

import "time"

type User struct {
	UserId       string    `json:"user_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
	CreatedAt    time.Time `json:"created_at" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at" validate:"required"`
	Email        string    `json:"email" example:"maxim@yandex.ru" validate:"required"`
	AvatarPath   *string   `json:"avatar_path"`
	PasswordHash string    `json:"password_hash" validate:"required"`
	IsActivated  bool      `json:"is_activated" example:"false" validate:"required"`
}

type CreatedUser struct {
	UserId         string `json:"user_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
	Email          string `json:"email" example:"maxim@yandex.ru" validate:"required"`
	ActivationLink string `json:"activation_link" validate:"required"`
}

type ChangePasswordCode struct {
	ChangePasswordCodeId int    `json:"change_password_code_id" db:"change_password_code_id" validate:"required"`
	Code                 string `json:"code" validate:"required"`
	UserId               string `json:"user_id" validate:"required"`
}
