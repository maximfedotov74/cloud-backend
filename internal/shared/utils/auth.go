package utils

import (
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
)

func Auth(r *http.Request) (*string, ex.Error) {
	user, ok := r.Context().Value(keys.UserSessionCtx).(string)

	if !ok {
		return nil, ex.NewErr(ex.Unauthorized, http.StatusUnauthorized)
	}

	return &user, nil

}
