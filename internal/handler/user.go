package handler

import (
	"log"
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/shared/utils"
)

type userService interface {
}

type UserHandler struct {
	userService userService
	router      *http.ServeMux
}

func NewUserHandler(userService userService, router *http.ServeMux) *UserHandler {
	return &UserHandler{userService: userService, router: router}
}

func (h *UserHandler) StartHandlers() {
	h.router.HandleFunc("GET /api/user/activate/{activationLink}", func(w http.ResponseWriter, r *http.Request) {
		ip, _ := utils.GetIP(r)
		log.Println("IP:", ip)
		http.Redirect(w, r, "https://ya.ru/", http.StatusPermanentRedirect)
	})
}
