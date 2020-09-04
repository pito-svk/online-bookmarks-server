package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

func NewAuthHandler(router *chi.Mux, us domain.AuthUsecase) {
	handler := &AuthHandler{
		AuthUsecase: us,
	}

	router.Post("/auth/register", handler.RegisterUser)
}

func (a *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := entity.User{}

	userResponse, err := a.AuthUsecase.Register(&user)

	if err != nil {
		w.Write([]byte("Error: TODO"))
	}

	// TODO: Implement
	w.Write([]byte(userResponse.Email))
}
