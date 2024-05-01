package model

import "time"

type Session struct {
	SessionId int       `json:"session_id" db:"session_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	UserId    string    `json:"user_id" db:"user_id" validate:"required"`
	UserAgent string    `json:"user_agent" db:"user_agent" validate:"required"`
	Token     string    `json:"-" db:"token" validate:"required"`
	Ip        string    `json:"-"`
}

type UserSessionsResponse struct {
	Current *Session  `json:"current" validate:"required"`
	All     []Session `json:"sessions" validate:"required"`
}

type LocalSession struct {
	UserId    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	Email     string `json:"email"`
	Ip        string `json:"-"`
}
