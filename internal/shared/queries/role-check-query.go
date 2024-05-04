package queries

import (
	"fmt"

	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
)

func GenerateRoleCheckerQuery(userId string, roles ...string) string {

	if len(roles) == 0 {
		roles = append(roles, keys.UserRole)
	}

	parameter := 1

	q := fmt.Sprintf(`
	SELECT COUNT(u.user_id)
	FROM %s u
	`, keys.UserTable)

	for _, role := range roles {
		a := fmt.Sprintf("ur%[1]d ON u.user_id = ur%[1]d.user_id", parameter)
		b := fmt.Sprintf("r%[1]d ON ur%[1]d.role_id = r%[1]d.role_id AND r%[1]d.title = '%s'", parameter, role)
		q += fmt.Sprintf(`
		INNER JOIN %s %s
		INNER JOIN %s %s
		`, keys.UserRoleTable, a, keys.RoleTable, b)
		parameter++
	}

	q += fmt.Sprintf(`
		WHERE u.user_id = '%s';
	`, userId)

	return q
}
