package mw

import (
	"context"
	"net/http"
	"strings"

	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/jwt"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
	"github.com/maximfedotov74/cloud-api/internal/shared/utils"
)

type sessionRepository interface {
	FindByAgentAndUserId(ctx context.Context, agent string, userId string) (*model.Session, ex.Error)
}

type userService interface {
	FindById(ctx context.Context, id string) (*model.User, ex.Error)
}

type AuthMW func(http.HandlerFunc) http.HandlerFunc

func NewAuthMW(userService userService, sessionRepository sessionRepository, jwtService *jwt.JwtService) AuthMW {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), keys.UserSessionCtx, nil)

			authHeader := r.Header.Get(keys.AuthorizationHeader)

			splittedHeader := strings.Split(authHeader, " ")

			authFall := ex.NewErr(ex.Unauthorized, http.StatusUnauthorized)

			if len(splittedHeader) != 2 {
				utils.WriteJSON(w, authFall.Status(), authFall)
				return
			}

			accessToken := splittedHeader[1]
			claims, err := jwtService.Parse(accessToken, jwt.AccessToken)
			if err != nil {
				utils.WriteJSON(w, authFall.Status(), authFall)
				return
			}
			user, fall := userService.FindById(r.Context(), claims.UserId)
			if fall != nil {
				utils.WriteJSON(w, fall.Status(), fall)
				return
			}
			session, fall := sessionRepository.FindByAgentAndUserId(r.Context(), claims.UserAgent, user.UserId)

			if fall != nil {
				accessToken, refreshToken := utils.RemoveTokensCookie()
				http.SetCookie(w, accessToken)
				http.SetCookie(w, refreshToken)
				forbidden := ex.NewErr(ex.Forbidden, http.StatusForbidden)
				utils.WriteJSON(w, forbidden.Status(), forbidden)
				return
			}

			localSession := model.LocalSession{UserId: session.UserId, UserAgent: session.UserAgent, Email: user.Email, Ip: session.Ip}
			ctx = context.WithValue(ctx, keys.UserSessionCtx, localSession)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
