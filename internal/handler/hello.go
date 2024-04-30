package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/maximfedotov74/cloud-api/internal/cfg"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/mw"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/utils"
)

var id int = 1

var msgs = []model.Hello{}

type HelloHandler struct {
	config *cfg.Config
	router *http.ServeMux
}

func NewHelloHandler(config *cfg.Config, router *http.ServeMux) *HelloHandler {
	return &HelloHandler{config: config, router: router}
}

func (h *HelloHandler) StartHandlers() {

	h.router.HandleFunc("POST /api/hello", h.createMessage)

	h.router.HandleFunc("GET /api/hello", h.getMessages)

	h.router.HandleFunc("GET /api/hello/{id}", h.getMessageById)

	h.router.HandleFunc("GET /api/auth", mw.AuthMw(h.auth))
}

// @Summary Create Message
// @Description Create Message
// @Tags message
// @Accept json
// @Produce json
// @Param dto body model.CreateMsgDto true "Create message with body dto"
// @Router /api/hello [post]
// @Success 201 {object} ex.AppErr
// @Failure 400 {object} ex.AppErr
// @Failure 404 {object} ex.AppErr
// @Failure 500 {object} ex.AppErr
func (h *HelloHandler) createMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, ex.ServerError(err.Error()))
		return
	}

	var dto model.CreateMsgDto

	err = json.Unmarshal(body, &dto)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, ex.ServerError(err.Error()))
		return
	}

	validate := validator.New()

	err = validate.Struct(&dto)

	if err != nil {
		error_messages := err.(validator.ValidationErrors)
		items := ex.ValidationMessages(error_messages)
		validError := ex.NewValidErr(items)
		utils.WriteJSON(w, http.StatusBadRequest, validError)
		return
	}

	msgs = append(msgs, model.Hello{
		Id:      id,
		Message: dto.Message,
	})
	id++

	utils.WriteJSON(w, http.StatusCreated, ex.GetCreated())
}

// @Summary Get all messages
// @Description Get all messages
// @Tags message
// @Accept json
// @Produce json
// @Router /api/hello [get]
// @Success 200 {array} model.Hello
// @Failure 400 {object} ex.AppErr
// @Failure 404 {object} ex.AppErr
// @Failure 500 {object} ex.AppErr
func (h *HelloHandler) getMessages(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, msgs)
}

// @Summary Get message by id
// @Description Get message by id
// @Tags message
// @Accept json
// @Produce json
// @Router /api/hello/{id} [get]
// @Param id path int true "message id"
// @Success 200 {object} model.Hello
// @Failure 400 {object} ex.AppErr
// @Failure 404 {object} ex.AppErr
// @Failure 500 {object} ex.AppErr
func (h *HelloHandler) getMessageById(w http.ResponseWriter, r *http.Request) {
	var msg *model.Hello = nil

	queryId := r.PathValue("id")

	id, err := strconv.Atoi(queryId)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, ex.NewErr(ex.VALIDATION_ID, http.StatusBadRequest))
		return
	}

	for _, item := range msgs {
		if item.Id == id {
			msg = &item
			break
		}
	}

	if msg == nil {
		utils.WriteJSON(w, http.StatusNotFound, ex.NewErr("Сообщение не найдено!", http.StatusNotFound))
		return
	}
	utils.WriteJSON(w, http.StatusOK, msg)
}

func (h *HelloHandler) auth(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(string)

	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, ex.NewErr(ex.UNAUTHORIZED, http.StatusUnauthorized))
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, model.User{Token: user})
}
