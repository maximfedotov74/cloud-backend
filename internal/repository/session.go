package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/msg"
	"github.com/maximfedotov74/cloud-api/internal/shared/db"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
)

type SessionRepository struct {
	db db.PostgresClient
}

func NewSessionRepository(db db.PostgresClient) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, dto dto.CreateSession) ex.Error {
	q := fmt.Sprintf("SELECT session_id, user_id, user_agent, refresh_token FROM %s WHERE user_id = $1 AND user_agent = $2;", keys.SessionTable)

	sessionModel := model.Session{}

	row := r.db.QueryRow(ctx, q, dto.UserId, dto.UserAgent)

	err := row.Scan(&sessionModel.SessionId, &sessionModel.UserId, &sessionModel.UserAgent, &sessionModel.Token)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			q = fmt.Sprintf("INSERT INTO %s (refresh_token, user_agent, user_id, ip) VALUES ($1, $2, $3, $4) RETURNING session_id;", keys.SessionTable)
			_, err = r.db.Exec(ctx, q, dto.Token, dto.UserAgent, dto.UserId, dto.Ip)
			if err != nil {
				log.Println(err.Error())
				return ex.ServerError(msg.SessionCreateError)
			}
			return nil
		}
		return ex.ServerError(err.Error())
	}
	q = fmt.Sprintf("UPDATE %s SET refresh_token = $1, updated_at = %s WHERE session_id = $2;", keys.SessionTable, keys.PsqlCurrentTimestamp)

	_, err = r.db.Exec(ctx, q, dto.Token, sessionModel.SessionId)

	if err != nil {
		return ex.ServerError(msg.SessionUpdateError)
	}

	return nil
}

func (r *SessionRepository) FindByAgentAndToken(ctx context.Context, agent string, token string) (*model.Session, ex.Error) {

	query := fmt.Sprintf("SELECT session_id, user_id, user_agent, refresh_token, created_at, updated_at, ip FROM %s WHERE user_agent = $1 AND token = $2;",
		keys.SessionTable)

	sessionModel := model.Session{}

	row := r.db.QueryRow(ctx, query, agent, token)

	err := row.Scan(&sessionModel.SessionId, &sessionModel.UserId, &sessionModel.UserAgent,
		&sessionModel.Token, &sessionModel.CreatedAt, &sessionModel.UpdatedAt, &sessionModel.Ip,
	)

	if err != nil {
		return nil, ex.NewErr(msg.SessionNotFound, http.StatusNotFound)
	}

	return &sessionModel, nil

}

func (r *SessionRepository) FindByAgentAndUserId(ctx context.Context, agent string, userId string) (*model.Session, ex.Error) {

	query := fmt.Sprintf("SELECT session_id, user_id, user_agent, refresh_token, created_at, updated_at, ip FROM %s WHERE user_agent = $1 AND user_id = $2;",
		keys.SessionTable)

	sessionModel := model.Session{}

	row := r.db.QueryRow(ctx, query, agent, userId)

	err := row.Scan(&sessionModel.SessionId, &sessionModel.UserId, &sessionModel.UserAgent,
		&sessionModel.Token, &sessionModel.CreatedAt, &sessionModel.UpdatedAt, &sessionModel.Ip,
	)

	if err != nil {
		log.Println(err.Error())
		return nil, ex.NewErr(msg.SessionNotFound, http.StatusNotFound)
	}

	return &sessionModel, nil

}

func (r *SessionRepository) RemoveSession(ctx context.Context, userId string, sessionId int) ex.Error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND session_id = $2;", keys.SessionTable)
	_, err := r.db.Exec(ctx, query, userId, sessionId)
	if err != nil {
		return ex.ServerError(err.Error())
	}
	return nil
}

func (r *SessionRepository) RemoveExceptCurrentSession(ctx context.Context, userId string, sessionId int) ex.Error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND session_id != $2;", keys.SessionTable)
	_, err := r.db.Exec(ctx, query, userId, sessionId)
	if err != nil {
		return ex.ServerError(err.Error())
	}
	return nil
}

func (r *SessionRepository) RemoveSessionByToken(ctx context.Context, token string) ex.Error {
	query := fmt.Sprintf("DELETE FROM %s WHERE refresh_token = $1;", keys.SessionTable)
	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return ex.ServerError(err.Error())
	}
	return nil
}

func (r *SessionRepository) RemoveAllSessions(ctx context.Context, userId string) ex.Error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1;", keys.SessionTable)

	_, err := r.db.Exec(ctx, query, userId)
	if err != nil {
		return ex.ServerError(err.Error())
	}
	return nil
}

func (r *SessionRepository) GetUserSessions(ctx context.Context, userId string, token string) (*model.UserSessionsResponse, ex.Error) {
	q := fmt.Sprintf(`SELECT session_id, user_id, user_agent, refresh_token, created_at, updated_at, ip
	FROM %s WHERE user_id = $1 ORDER BY updated_at DESC;`, keys.SessionTable)

	var sessions []model.Session

	rows, err := r.db.Query(ctx, q, userId)

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}

	response := model.UserSessionsResponse{}

	for rows.Next() {
		s := model.Session{}

		err := rows.Scan(&s.SessionId, &s.UserId, &s.UserAgent,
			&s.Token, &s.CreatedAt, &s.UpdatedAt, &s.Ip,
		)

		if err != nil {
			return nil, ex.NewErr(msg.SessionNotFound, http.StatusNotFound)
		}

		if s.Token == token {
			response.Current = &s
		}

		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, ex.ServerError(err.Error())
	}

	response.All = sessions
	return &response, nil
}
