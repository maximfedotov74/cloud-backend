package model

import "github.com/maximfedotov74/cloud-api/internal/shared/jwt"

type LoginResponse struct {
	UserId string     `json:"user_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
	Tokens jwt.Tokens `json:"tokens" validate:"required"`
}

type RegistrationResponse struct {
	UserId string `json:"user_id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" validate:"required"`
}
