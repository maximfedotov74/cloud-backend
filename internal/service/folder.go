package service

import (
	"context"
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/msg"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
)

type folderRepository interface {
	Create(ctx context.Context, title string, parentId *string, userId string) ex.Error
	FindByTitleInFolder(ctx context.Context, parentFolderId string, title string) (bool, ex.Error)
	FindById(ctx context.Context, folderId string) (*model.Folder, ex.Error)
}

type FolderService struct {
	folderRepository folderRepository
}

func NewFolderService(folderRepository folderRepository) *FolderService {
	return &FolderService{folderRepository: folderRepository}
}

func (s *FolderService) Create(ctx context.Context, title string, parentId *string, userId string) ex.Error {
	var parentFolderId *string

	if parentId != nil {
		p, fall := s.folderRepository.FindById(ctx, *parentId)
		if fall != nil {
			return fall
		}

		existsWithName, fall := s.folderRepository.FindByTitleInFolder(ctx, p.FolderId, title)

		if fall != nil {
			return fall
		}

		if existsWithName {
			return ex.NewErr(msg.FolderNameExists, http.StatusBadRequest)
		}

		parentFolderId = &p.FolderId
	}
	return s.folderRepository.Create(ctx, title, parentFolderId, userId)
}
