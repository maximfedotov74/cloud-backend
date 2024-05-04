package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/msg"
	"github.com/maximfedotov74/cloud-api/internal/shared/db"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
)

type FileRepository struct {
	db db.PostgresClient
}

func NewFileRepository(db db.PostgresClient) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) FindById(ctx context.Context, fileId string) (*model.File, ex.Error) {
	q := fmt.Sprintf(`SELECT file_id, created_at, updated_at, deleted_at, title, ext, size, folder_id, user_id FROM
  %s WHERE file_id = $1 AND deleted_at IS NULL;`, keys.FileTable)

	f := model.File{}

	row := r.db.QueryRow(ctx, q, fileId)

	err := row.Scan(&f.FileId, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt, &f.Title, &f.Ext, &f.Size, &f.FolderId, &f.UserId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ex.NewErr(msg.FileNotFound, http.StatusNotFound)
		}
		return nil, ex.ServerError(err.Error())
	}
	return &f, nil
}

func (r *FileRepository) FindByTitleInFolder(ctx context.Context, folderId string, title string) (bool, ex.Error) {
	q := fmt.Sprintf(`SELECT file_id FROM %s
  WHERE folder_id = $1 AND CONCAT(title,ext) = $2 AND deleted_at IS NULL`, keys.FileTable)
	var fileId string

	row := r.db.QueryRow(ctx, q, folderId, title)

	err := row.Scan(&fileId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return true, ex.ServerError(err.Error())
	}

	if fileId == "" {
		return false, nil
	}
	return true, nil
}

func (r *FileRepository) Create(ctx context.Context, input dto.CreateFile) (*string, ex.Error) {

	q := fmt.Sprintf("INSERT INTO %s (title, ext, size, folder_id, user_id) VALUES ($1,$2,$3,$4,$5) RETURNING file_id;", keys.FileTable)

	row := r.db.QueryRow(ctx, q, input.Title, input.Ext, input.Size, input.FolderId, input.UserId)

	var fileId string

	err := row.Scan(&fileId)

	if err != nil {
		return nil, ex.ServerError(msg.FileCreateError)
	}

	return &fileId, nil
}
