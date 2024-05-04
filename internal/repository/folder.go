package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/msg"
	"github.com/maximfedotov74/cloud-api/internal/shared/db"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/keys"
)

type FolderRepository struct {
	db db.PostgresClient
}

func NewFolderRepository(db db.PostgresClient) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) FindById(ctx context.Context, folderId string) (*model.Folder, ex.Error) {
	q := fmt.Sprintf(`SELECT folder_id, created_at, updated_at, deleted_at, title, parent_folder_id, user_id FROM %s
  WHERE folder_id = $1 AND deleted_at IS NULL;`, keys.FolderTable)

	f := model.Folder{}

	row := r.db.QueryRow(ctx, q, folderId)

	err := row.Scan(&f.FolderId, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt, &f.Title, &f.ParentFolderId, &f.UserId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ex.NewErr(msg.FolderNotFound, http.StatusNotFound)
		}
		return nil, ex.ServerError(err.Error())
	}
	return &f, nil
}

func (r *FolderRepository) FindByTitleInFolder(ctx context.Context, parentFolderId string, title string) (bool, ex.Error) {
	q := fmt.Sprintf(`SELECT folder_id FROM %s
  WHERE parent_folder_id = $1 AND title = $2 AND deleted_at IS NULL`, keys.FolderTable)

	var folderId string

	row := r.db.QueryRow(ctx, q, parentFolderId, title)

	err := row.Scan(&folderId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return true, ex.ServerError(err.Error())
	}

	if folderId == "" {
		return false, nil
	}

	return true, nil
}

func (r *FolderRepository) Create(ctx context.Context, title string, parentId *string, userId string) ex.Error {

	q := fmt.Sprintf("INSERT INTO %s (title,parent_folder_id,user_id) VALUES ($1,$2,$3);", keys.FolderTable)

	_, err := r.db.Exec(ctx, q, title, parentId, userId)

	if err != nil {
		return ex.ServerError(msg.FolderCreateError)
	}

	return nil
}
