package repository

import (
	"context"
	"fmt"
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/msg"
	"github.com/maximfedotov74/cloud-api/internal/shared/db"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
	"github.com/maximfedotov74/cloud-api/internal/shared/queries"
)

const (
// addRoleToUser =
)

type RoleRepository struct {
	db db.PostgresClient
}

func NewRoleRepository(db db.PostgresClient) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) Create(ctx context.Context, input dto.CreateRole) (*model.Role, ex.Error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING role_id, title;", keys.RoleTable)
	row := r.db.QueryRow(ctx, query, input.Title)
	role := model.Role{}
	err := row.Scan(&role.Id, &role.Title)

	if err != nil {
		return nil, ex.NewErr(msg.RoleCreateError, http.StatusInternalServerError)
	}
	return &role, nil
}

func (r *RoleRepository) FindAll(ctx context.Context) ([]model.UserRole, ex.Error) {
	q := fmt.Sprintf("SELECT role_id, title FROM %s;", keys.RoleTable)

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}

	var roles []model.UserRole

	for rows.Next() {
		r := model.UserRole{}

		err := rows.Scan(&r.Id, &r.Title)
		if err != nil {
			return nil, ex.ServerError(err.Error())
		}
		roles = append(roles, r)
	}

	if err := rows.Err(); err != nil {
		return nil, ex.ServerError(err.Error())
	}

	return roles, nil

}

func (r *RoleRepository) FindByTitle(ctx context.Context, title string) (*model.Role, ex.Error) {
	query := fmt.Sprintf(`
	SELECT r.role_id, r.title, u.user_id, u.email FROM %s as r
	LEFT JOIN %s as ur ON r.role_id = ur.role_id
	LEFT JOIN %s as u ON u.user_id = ur.user_id
	WHERE r.title = $1;`, keys.RoleTable, keys.UserRoleTable, keys.UserTable)
	rows, err := r.db.Query(ctx, query, title)

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}
	defer rows.Close()
	founded := false
	role := model.Role{}

	users := []model.RoleUser{}

	for rows.Next() {
		user := model.RoleUser{}
		err := rows.Scan(&role.Id, &role.Title, &user.Id, &user.Email)
		if err != nil {
			return nil, ex.ServerError(err.Error())
		}
		if user.Id != nil {
			users = append(users, user)
		}
		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, ex.ServerError(err.Error())
	}

	if !founded {
		return nil, ex.NewErr(msg.RoleNotFound, http.StatusNotFound)
	}

	role.Users = users

	return &role, nil

}

func (r *RoleRepository) FindWithUsers(ctx context.Context) ([]model.Role, ex.Error) {

	query := fmt.Sprintf(`
	SELECT r.role_id, r.title,  u.user_id, u.email, ur.user_role_id, ur.role_id as ur_role_id FROM %s as r
	LEFT JOIN %s as ur ON r.role_id = ur.role_id
	LEFT JOIN %s as u ON u.user_id = ur.user_id;`, keys.RoleTable, keys.UserRoleTable, keys.UserTable)

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}
	defer rows.Close()

	founded := false
	rolesMap := make(map[int]model.Role)
	usersMap := make(map[int]model.RoleUser)

	var rolesOrder []int
	var usersOrder []int

	for rows.Next() {
		var role model.Role
		role.Users = []model.RoleUser{}
		var user model.RoleUser

		err := rows.Scan(&role.Id, &role.Title, &user.Id, &user.Email, &user.UserRoleId, &user.RoleId)
		if err != nil {
			return nil, ex.ServerError(err.Error())
		}

		if user.UserRoleId != nil {
			_, ok := usersMap[*user.UserRoleId]
			if !ok {
				usersMap[*user.UserRoleId] = user
				usersOrder = append(usersOrder, *user.UserRoleId)
			}
		}
		_, ok := rolesMap[role.Id]
		if !ok {
			rolesMap[role.Id] = role
			rolesOrder = append(rolesOrder, role.Id)
		}
		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, ex.ServerError(err.Error())
	}

	roles := make([]model.Role, 0, len(rolesMap))

	if !founded {
		return roles, nil
	}

	for _, key := range usersOrder {
		user := usersMap[key]
		role := rolesMap[*user.RoleId]
		role.Users = append(role.Users, user)
		rolesMap[role.Id] = role
	}

	for _, key := range rolesOrder {
		role := rolesMap[key]
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *RoleRepository) CheckRolesInUser(ctx context.Context, userId string, roles ...string) bool {
	q := queries.GenerateRoleCheckerQuery(userId, roles...)

	row := r.db.QueryRow(ctx, q)

	count := 0

	err := row.Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func (r *RoleRepository) AddToUser(ctx context.Context, roleId int, userId string, tx db.Transaction) ex.Error {

	q := fmt.Sprintf("INSERT INTO %s (user_id, role_id) VALUES ($1, $2);", keys.UserRoleTable)

	if tx != nil {
		_, err := tx.Exec(ctx, q, userId, roleId)
		if err != nil {
			return ex.NewErr(msg.RoleAddError, http.StatusInternalServerError)
		}
		return nil
	}
	_, err := r.db.Exec(ctx, q, userId, roleId)
	if err != nil {
		return ex.NewErr(msg.RoleAddError, http.StatusInternalServerError)
	}
	return nil
}

func (r *RoleRepository) Remove(ctx context.Context, roleId int) ex.Error {
	query := fmt.Sprintf("DELETE FROM %s WHERE role_id = $1;", keys.RoleTable)

	_, err := r.db.Exec(ctx, query, roleId)

	if err != nil {
		return ex.ServerError(err.Error())
	}
	return nil
}

func (r *RoleRepository) RemoveFromUser(ctx context.Context, roleId int, userId string) ex.Error {

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND role_id = $2;", keys.UserRoleTable)
	_, err := r.db.Exec(ctx, query, userId, roleId)

	if err != nil {
		return ex.NewErr(msg.RoleDeleteError, http.StatusInternalServerError)
	}

	return nil
}
