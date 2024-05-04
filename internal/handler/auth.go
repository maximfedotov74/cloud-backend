package handler

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/utils"
)

type authService interface {
	Login(ctx context.Context, input dto.CreateUser, userAgent string, ip net.IP) (*model.LoginResponse, ex.Error)
	Registration(ctx context.Context, input dto.CreateUser) (*string, ex.Error)
	Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, ex.Error)
	Logout(ctx context.Context, token string) ex.Error
}

type AuthHandler struct {
	authService authService
	router      *http.ServeMux
}

func NewAuthHandler(authService authService, router *http.ServeMux) *AuthHandler {
	return &AuthHandler{authService: authService, router: router}
}

func (h *AuthHandler) StartHandlers() {
	h.router.HandleFunc("POST /api/auth/registration", h.registration)
	h.router.HandleFunc("POST /api/auth/login", h.login)
	h.router.HandleFunc("POST /api/auth/logout", h.logout)
	h.router.HandleFunc("GET /api/auth/refresh-token", h.refreshTokens)
}

// @Summary Registration is system
// @Description Registration is system
// @Tags auth
// @Accept json
// @Produce json
// @Param dto body dto.CreateUser true "Registration is system with body dto"
// @Router /api/auth/registration [post]
// @Success 201 {object} model.RegistrationResponse
// @Failure 400 {object} ex.AppErr
// @Failure 404 {object} ex.AppErr
// @Failure 500 {object} ex.AppErr
func (h *AuthHandler) registration(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, ex.ServerError(err.Error()))
		return
	}

	var input dto.CreateUser

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

	userId, fall := h.authService.Registration(r.Context(), input)

	if fall != nil {
		utils.WriteJSON(w, fall.Status(), fall)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, model.RegistrationResponse{UserId: *userId})
}

// @Summary Login
// @Description Login to an account with account data
// @Tags auth
// @Accept json
// @Produce json
// @Param dto body dto.CreateUser true "login in account"
// @Router /api/auth/login [post]
// @Success 201 {object} model.LoginResponse
// @Failure 400 {object} ex.ValidationError
// @Failure 404 {object} ex.AppErr
// @Failure 500 {object} ex.AppErr
func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, ex.ServerError(err.Error()))
		return
	}

	var input dto.CreateUser

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

	userAgent := r.UserAgent()
	ip, err := utils.GetIP(r)
	if err != nil {
		fall := ex.ServerError(err.Error())
		utils.WriteJSON(w, fall.Status(), fall)
		return
	}

	res, fall := h.authService.Login(r.Context(), input, userAgent, ip)
	if fall != nil {
		utils.WriteJSON(w, fall.Status(), fall)
		return
	}
	accessToken, refreshToken := utils.SetTokensCookie(res.Tokens)
	http.SetCookie(w, accessToken)
	http.SetCookie(w, refreshToken)
	utils.WriteJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {}

func (h *AuthHandler) refreshTokens(w http.ResponseWriter, r *http.Request) {}
