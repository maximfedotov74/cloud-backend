package service

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/msg"
	"github.com/maximfedotov74/cloud-api/internal/shared/db"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/file"
)

type fileRepository interface {
	Create(ctx context.Context, input dto.CreateFile) (*string, ex.Error)
	FindByTitleInFolder(ctx context.Context, folderId string, title string) (bool, ex.Error)
	FindById(ctx context.Context, fileId string) (*model.File, ex.Error)
}

type fileFolderRepository interface {
	FindById(ctx context.Context, folderId string) (*model.Folder, ex.Error)
}

type FileService struct {
	fileRepository     fileRepository
	folderRepository   fileFolderRepository
	transactionManager db.TransactionManager
	fileStorage        *file.FileClient
}

func NewFileService(fileRepository fileRepository, folderRepository fileFolderRepository, transactionManager db.TransactionManager,
	fileStorage *file.FileClient,
) *FileService {
	return &FileService{fileRepository: fileRepository, folderRepository: folderRepository, transactionManager: transactionManager, fileStorage: fileStorage}
}

func (s *FileService) Create(ctx context.Context, input dto.CreateFile, h *multipart.FileHeader) ex.Error {

	tx, err := s.transactionManager.Begin(ctx)
	var fall ex.Error = nil

	if err != nil {
		fall = ex.ServerError(err.Error())
		return fall
	}

	defer func() {
		if fall != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	p, fall := s.folderRepository.FindById(ctx, input.FolderId)
	if fall != nil {
		return fall
	}

	existsWithName, fall := s.fileRepository.FindByTitleInFolder(ctx, p.FolderId, input.Title+input.Ext)

	if fall != nil {
		return fall
	}

	if existsWithName {
		return ex.NewErr(msg.FileNameExists, http.StatusBadRequest)
	}
	fileId, fall := s.fileRepository.Create(ctx, input)

	if fall != nil {
		return fall
	}

	_, err = s.fileStorage.Upload(ctx, input.UserId, *fileId, h)

	if err != nil {
		fall = ex.ServerError(msg.FileSaveToStorageError)
		return fall
	}
	return nil
}
