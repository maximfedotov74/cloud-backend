package utils

import (
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/shared/jwt"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
)

func SetTokensCookie(tokens jwt.Tokens) (*http.Cookie, *http.Cookie) {
	accessToken := &http.Cookie{
		Name:    keys.AccessToken,
		Value:   tokens.AccessToken,
		Expires: tokens.AccessExpTime,
	}
	refreshToken := &http.Cookie{
		Name:     keys.RefreshToken,
		Value:    tokens.RefreshToken,
		Expires:  tokens.RefreshExpTime,
		HttpOnly: true,
	}
	return accessToken, refreshToken
}

func RemoveTokensCookie() (*http.Cookie, *http.Cookie) {
	accessToken := &http.Cookie{
		Name:   keys.AccessToken,
		Value:  "",
		MaxAge: -1,
	}
	refreshToken := &http.Cookie{
		Name:     keys.RefreshToken,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}
	return accessToken, refreshToken
}
