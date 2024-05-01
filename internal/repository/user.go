package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/msg"
	"github.com/maximfedotov74/cloud-api/internal/shared/db"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
)

type userRoleRepository interface {
	FindByTitle(ctx context.Context, title string) (*model.Role, ex.Error)
	AddToUser(ctx context.Context, roleId int, userId string, tx db.Transaction) ex.Error
}

type UserRepository struct {
	db             db.PostgresClient
	roleRepository userRoleRepository
}

func NewUserRepository(db db.PostgresClient, roleRepository userRoleRepository) *UserRepository {
	return &UserRepository{db: db, roleRepository: roleRepository}
}

func (r *UserRepository) Create(ctx context.Context, input dto.CreateUser, tx db.Transaction) (*model.CreatedUser, ex.Error) {

	q := fmt.Sprintf("INSERT INTO %s (email, password_hash) VALUES ($1,$2) RETURNING user_id, email;", keys.UserTable)

	row := tx.QueryRow(ctx, q, input.Email, input.Password)

	var email string
	var id string
	err := row.Scan(&id, &email)

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}

	role, fall := r.roleRepository.FindByTitle(ctx, keys.UserRole)

	if fall != nil {
		return nil, fall
	}

	fall = r.roleRepository.AddToUser(ctx, role.Id, id, tx)

	if fall != nil {
		return nil, fall
	}
	q = fmt.Sprintf("INSERT INTO %s (user_id, activation_account_link) VALUES ($1, %s) RETURNING activation_account_link;",
		keys.UserActivationTable, keys.PsqlUUID)

	row = tx.QueryRow(ctx, q, id)

	var link string

	err = row.Scan(&link)

	if err != nil {
		fall = ex.ServerError(err.Error())
		return nil, fall
	}

	return &model.CreatedUser{UserId: id, Email: email, ActivationLink: link}, nil
}

func (r *UserRepository) Update(ctx context.Context, input dto.UpdateUser, id int) ex.Error {

	var queries []string

	if input.AvatarPath != nil {
		queries = append(queries, fmt.Sprintf("avatar_path = '%s'", *input.AvatarPath))
	}

	queries = append(queries, fmt.Sprintf("updated_at = %s", keys.PsqlCurrentTimestamp))

	if len(queries) > 0 {
		q := fmt.Sprintf("UPDATE %s SET %s WHERE user_id = $1", keys.UserTable, strings.Join(queries, ","))
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return ex.ServerError(fmt.Sprintf("%s, details: \n %s", msg.UpdateUserError, err.Error()))
		}
		return nil
	}

	return nil
}

func (r *UserRepository) findByField(ctx context.Context, field string, value any) (*model.User, ex.Error) {
	q := fmt.Sprintf(`SELECT user_id, created_at, updated_at, email, avatar_path, password_hash, is_activated FROM
	%s WHERE %s = $1;`, keys.UserTable, field)

	row := r.db.QueryRow(ctx, q, value)

	u := model.User{}

	err := row.Scan(&u.UserId, &u.CreatedAt, &u.UpdatedAt, &u.Email, &u.AvatarPath, &u.PasswordHash, &u.IsActivated)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ex.NewErr(msg.UserNotFound, http.StatusNotFound)
		}
		return nil, ex.ServerError(err.Error())
	}

	return &u, nil
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*model.User, ex.Error) {
	return r.findByField(ctx, "user_id", id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, ex.Error) {
	return r.findByField(ctx, "email", email)
}

func (r *UserRepository) RemoveChangePasswordCode(ctx context.Context, userId string, tx db.Transaction) error {
	query := "DELETE FROM change_password_code WHERE user_id = $1"

	if tx != nil {
		_, err := tx.Exec(ctx, query, userId)

		if err != nil {
			return err
		}
		return nil
	}

	_, err := r.db.Exec(ctx, query, userId)

	if err != nil {
		return err
	}
	return nil

}

func (r *UserRepository) FindChangePasswordCode(ctx context.Context, userId string, code string) (*model.ChangePasswordCode, ex.Error) {
	query := fmt.Sprintf(`SELECT change_password_code_id, code, user_id FROM
	%s WHERE user_id = $1 AND code = $2 AND end_time > $3;`, keys.UserChangePasswordCodeTable)

	row := r.db.QueryRow(ctx, query, userId, code, keys.PsqlCurrentTimestamp)

	codeModel := model.ChangePasswordCode{}

	err := row.Scan(&codeModel.ChangePasswordCodeId, &codeModel.Code, &codeModel.UserId)

	if err != nil {
		return nil, ex.NewErr(msg.ChangePasswordCodeNotFound, http.StatusNotFound)
	}

	return &codeModel, nil

}

func (r *UserRepository) CreateChangePasswordCode(ctx context.Context, userId string) (*string, ex.Error) {

	var fall ex.Error = nil

	tx, err := r.db.Begin(ctx)

	if err != nil {
		fall = ex.ServerError(err.Error())
		return nil, fall
	}

	defer func() {
		if fall != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = r.RemoveChangePasswordCode(ctx, userId, tx)

	if err != nil {
		fall = ex.ServerError(err.Error())
		return nil, fall
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id) VALUES ($1) RETURNING code;", keys.UserChangePasswordCodeTable)

	row := tx.QueryRow(ctx, query, userId)

	var code string

	err = row.Scan(&code)

	if err != nil {
		fall = ex.ServerError(msg.CreateChangeCodeError)
		return nil, fall
	}

	return &code, nil
}

func (r *UserRepository) ChangePassword(ctx context.Context, userId string, newPassword string) ex.Error {

	query := fmt.Sprintf(`UPDATE %s SET password_hash = $1,
	updated_at = $2 WHERE user_id = $3;`, keys.UserTable)

	_, err := r.db.Exec(ctx, query, newPassword, keys.PsqlCurrentTimestamp, userId)
	if err != nil {
		return ex.ServerError(msg.UpdatePasswordError)
	}

	return nil
}
