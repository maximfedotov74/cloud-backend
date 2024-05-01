package mw

import (
	"context"
	"net/http"
	"strings"

	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
	"github.com/maximfedotov74/cloud-api/internal/shared/utils"
)

func AuthMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), keys.UserSessionCtx, nil)

		authHeader := r.Header.Get("Authorization")

		splittedHeader := strings.Split(authHeader, " ")

		if len(splittedHeader) != 2 {
			utils.WriteJSON(w, http.StatusUnauthorized, ex.NewErr(ex.Unauthorized, http.StatusUnauthorized))
			return
		}

		accessToken := splittedHeader[1]
		ctx = context.WithValue(ctx, keys.UserSessionCtx, accessToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
