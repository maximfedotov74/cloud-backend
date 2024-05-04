package dto

import "net"

type CreateSession struct {
	UserId    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	Token     string `json:"token"`
	Ip        net.IP `json:"-"`
}
