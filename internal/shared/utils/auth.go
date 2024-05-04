package utils

import (
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
)

func GetLocalSession(r *http.Request) (*model.LocalSession, ex.Error) {
	user, ok := r.Context().Value(keys.UserSessionCtx).(model.LocalSession)
	if !ok {
		return nil, ex.NewErr(ex.Unauthorized, http.StatusUnauthorized)
	}
	return &user, nil
}
