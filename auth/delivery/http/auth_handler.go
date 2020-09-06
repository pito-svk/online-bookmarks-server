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
	JwtSecret   string
}

func NewAuthHandler(router *chi.Mux, usecase domain.AuthUsecase, logger *logrus.Logger, jwtSecret string) {
	handler := &AuthHandler{
		AuthUsecase: usecase,
		Logger:      logger,
		JwtSecret:   jwtSecret,
	}

	router.Post("/auth/register", handler.RegisterUser)
}

type UserDataInput struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type AuthData struct {
	Token string `json:"token"`
}

type userCreatedResponse struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	AuthData  AuthData `json:"authData"`
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

func deliverErrorParsingJSONBodyError(w http.ResponseWriter) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: "Error parsing JSON body"})
}

func deliverBadRequestError(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: err.Error()})
}

func (a *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userData := UserDataInput{}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		a.Logger.Error(err)
		deliverErrorParsingJSONBodyError(w)
		return
	}

	err = validateCreateUserInput(&userData)
	if err != nil {
		deliverBadRequestError(w, err)
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

	authToken, err := a.AuthUsecase.GenerateAuthToken(userResponse.ID, a.JwtSecret)

	authData := AuthData{
		Token: authToken,
	}

	response := userCreatedResponse{
		ID:        userResponse.ID,
		Email:     userResponse.Email,
		FirstName: userResponse.FirstName,
		LastName:  userResponse.LastName,
		AuthData:  authData,
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response)
	return
}
