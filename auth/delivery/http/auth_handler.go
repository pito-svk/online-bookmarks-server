package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
	Logger      *logrus.Logger
}

func NewAuthHandler(router *chi.Mux, us domain.AuthUsecase, logger *logrus.Logger) {
	handler := &AuthHandler{
		AuthUsecase: us,
		Logger:      logger,
	}

	router.Post("/auth/register", handler.RegisterUser)
}

type UserDataInput struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type userCreatedResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type httpErrorMessage struct {
	Error string `json:"error"`
}

func validateCreateUserInput(userData *UserDataInput) error {
	v := validator.New()

	err := v.Struct(userData)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		for _, err := range validationErrors {
			if err.Tag() == "required" {
				err.StructField()
				return fmt.Errorf("Missing %s%s", strings.ToLower(string(err.Field()[0])), err.Field()[1:])
			}

			if err.Tag() == "email" {
				return errors.New("Invalid email")
			}
		}
	}

	return err
}

func (a *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userData := UserDataInput{}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		w.WriteHeader(400)

		a.Logger.Error(err)

		json.NewEncoder(w).Encode(httpErrorMessage{Error: "Error parsing JSON body"})
		return
	}

	err = validateCreateUserInput(&userData)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(httpErrorMessage{Error: err.Error()})
		return
	}

	userObject := entity.User{
		Email:     userData.Email,
		Password:  userData.Password,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
	}

	userResponse, err := a.AuthUsecase.Register(&userObject)
	if err != nil {
		w.WriteHeader(409)
		json.NewEncoder(w).Encode(httpErrorMessage{Error: err.Error()})
		return
	}

	response := userCreatedResponse{
		ID:        userResponse.ID,
		Email:     userResponse.Email,
		FirstName: userResponse.FirstName,
		LastName:  userResponse.LastName,
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response)
	return
}
