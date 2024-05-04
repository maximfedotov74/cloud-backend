package handler

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/mw"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/utils"
)

type fileService interface {
	Create(ctx context.Context, input dto.CreateFile, h *multipart.FileHeader) ex.Error
}

type FileHandler struct {
	fileService fileService
	router      *http.ServeMux
	authMW      mw.AuthMW
}

func NewFileHandler(fileService fileService, router *http.ServeMux, authMW mw.AuthMW) *FileHandler {
	return &FileHandler{fileService: fileService, router: router, authMW: authMW}
}

func (h *FileHandler) StartHandlers() {
	h.router.HandleFunc("POST /api/file", h.authMW(h.create))
}

// @Summary Create file
// @Description Create file
// @Tags file
// @Accept json
// @Produce json
// @Param folderId query string true "folder id"
// @Param file formData file true "File"
// @Router /api/file [post]
// @Success 201 {object} ex.AppErr
// @Failure 400 {object} ex.AppErr
// @Failure 404 {object} ex.AppErr
// @Failure 500 {object} ex.AppErr
func (h *FileHandler) create(w http.ResponseWriter, r *http.Request) {
	localSession, fall := utils.GetLocalSession(r)

	if fall != nil {
		utils.WriteJSON(w, fall.Status(), fall)
		return
	}

	_, header, err := r.FormFile("file")
	if err != nil {
		fall = ex.ServerError("Error when get access file")
		utils.WriteJSON(w, fall.Status(), fall)
		return
	}

	folderId := r.URL.Query().Get("folderId")

	if folderId == "" {
		fall = ex.NewErr(ex.ValidationId, http.StatusBadRequest)
		utils.WriteJSON(w, fall.Status(), fall)
		return
	}

	input := dto.CreateFile{Title: utils.GetFileName(header.Filename), Ext: utils.GetFileExt(header.Filename), Size: header.Size, FolderId: folderId, UserId: localSession.UserId}

	fall = h.fileService.Create(r.Context(), input, header)

	if fall != nil {
		utils.WriteJSON(w, fall.Status(), fall)
		return
	}

	created := ex.GetCreated()
	utils.WriteJSON(w, created.Status(), created)

}
