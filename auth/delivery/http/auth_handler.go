package http

import (
	"encoding/json"
	"errors"
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

type userDataInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func validateCreateUserInput(userData *userDataInput) error {
	if userData.Email == "" {
		return errors.New("Missing email")
	}

	if userData.Password == "" {
		return errors.New("Missing password")
	}

	if userData.FirstName == "" {
		return errors.New("Missing firstName")
	}

	if userData.LastName == "" {
		return errors.New("Mssing lastName")
	}

	return nil
}

func (a *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userData := userDataInput{}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = validateCreateUserInput(&userData)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	user := entity.User{
		Email:     userData.Email,
		Password:  userData.Password,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
	}

	userResponse, err := a.AuthUsecase.Register(&user)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Implement
	w.Write([]byte(userResponse.Email))
	return
}
