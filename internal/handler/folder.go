package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/mw"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/utils"
)

type folderService interface {
	Create(ctx context.Context, title string, parentId *string, userId string) ex.Error
}

type FolderHandler struct {
	folderService folderService
	router        *http.ServeMux
	authMW        mw.AuthMW
}

func NewFolderHandler(folderService folderService, router *http.ServeMux, authMw mw.AuthMW) *FolderHandler {
	return &FolderHandler{folderService: folderService, router: router, authMW: authMw}
}

func (h *FolderHandler) StartHandlers() {
	h.router.HandleFunc("POST /api/folder", h.authMW(h.create))
}

// @Summary Create folder
// @Description Create folder
// @Tags folder
// @Accept json
// @Produce json
// @Param dto body dto.CreateFolder true "Create folder with body dto"
// @Router /api/folder [post]
// @Success 201 {object} ex.AppErr
// @Failure 400 {object} ex.AppErr
// @Failure 404 {object} ex.AppErr
// @Failure 500 {object} ex.AppErr
func (h *FolderHandler) create(w http.ResponseWriter, r *http.Request) {

	localSession, fall := utils.GetLocalSession(r)

	if fall != nil {
		utils.WriteJSON(w, fall.Status(), fall)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, ex.ServerError(err.Error()))
		return
	}

	var input dto.CreateFolder

	err = json.Unmarshal(body, &input)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, ex.ServerError(err.Error()))
		return
	}

	validate := validator.New()

	err = validate.Struct(&input)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := ex.ValidationMessages(error_messages)
		validError := ex.NewValidErr(items)
		utils.WriteJSON(w, http.StatusBadRequest, validError)
		return
	}

	fall = h.folderService.Create(r.Context(), input.Title, input.ParentFolderId, localSession.UserId)
	if fall != nil {
		utils.WriteJSON(w, fall.Status(), fall)
		return
	}
	created := ex.GetCreated()
	utils.WriteJSON(w, created.Status(), created)

}
