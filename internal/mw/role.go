package mw

import (
	"context"
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/utils"
)

type roleRepository interface {
	CheckRolesInUser(ctx context.Context, userId string, roles ...string) bool
}

type RoleMw func(...string) func(http.HandlerFunc) http.HandlerFunc

func NewRoleMW(roleRepository roleRepository) RoleMw {
	return func(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				localSession, fall := utils.GetLocalSession(r)
				if fall != nil {
					utils.WriteJSON(w, fall.Status(), fall)
					return
				}
				flag := roleRepository.CheckRolesInUser(r.Context(), localSession.UserId, roles...)
				if !flag {
					forbidden := ex.NewErr(ex.Forbidden, http.StatusForbidden)
					utils.WriteJSON(w, forbidden.Status(), forbidden)
					return
				}
				next.ServeHTTP(w, r)
			}
		}
	}
}
