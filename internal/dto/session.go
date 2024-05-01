package dto

type CreateSession struct {
	UserId    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	Token     string `json:"token"`
	Ip        string `json:"-"`
}
