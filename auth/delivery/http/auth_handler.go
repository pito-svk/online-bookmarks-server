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

type userDataInput struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type authData struct {
	Token string `json:"token"`
}

type userCreatedResponse struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	AuthData  authData `json:"authData"`
}

type httpErrorMessage struct {
	Error string `json:"error"`
}

func setJSONContentTypeInResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func parseUserDataFromRequestBody(r *http.Request) (*userDataInput, error) {
	userData := userDataInput{}

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	return &userData, nil
}

func validateCreateUserInput(userData *userDataInput) error {
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

func deliverErrorParsingJSONBodyHttpError(w http.ResponseWriter) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: "Error parsing JSON body"})
}

func deliverBadRequestHttpError(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: err.Error()})
}

func deliverConflictHttpError(w http.ResponseWriter, err error) {
	w.WriteHeader(409)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: err.Error()})
}

func deliverInternalServerErrorHttpError(w http.ResponseWriter) {
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: "Internal Server Error"})
}

func deliverUserCreatedResponse(w http.ResponseWriter, response userCreatedResponse) {
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response)
}

func composeUserObjectFromUserData(userData *userDataInput) entity.User {
	return entity.User{
		Email:     userData.Email,
		Password:  userData.Password,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
	}
}

func composeUserCreatedResponse(user *entity.User, authToken string) userCreatedResponse {
	auth := authData{
		Token: authToken,
	}

	return userCreatedResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AuthData:  auth,
	}
}

func (authH *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	setJSONContentTypeInResponse(w)

	userData, err := parseUserDataFromRequestBody(r)
	if err != nil {
		authH.Logger.Error(err)
		deliverErrorParsingJSONBodyHttpError(w)
		return
	}
	err = validateCreateUserInput(userData)
	if err != nil {
		deliverBadRequestHttpError(w, err)
		return
	}

	userObject := composeUserObjectFromUserData(userData)
	userResponse, err := authH.AuthUsecase.Register(&userObject)
	if err != nil {
		if err.Error() == "User already exists" {
			deliverConflictHttpError(w, err)
		} else {
			authH.Logger.Error(err)
			deliverInternalServerErrorHttpError(w)
		}
		return
	}
	authToken, err := authH.AuthUsecase.GenerateAuthToken(userResponse.ID, authH.JwtSecret)
	response := composeUserCreatedResponse(userResponse, authToken)

	deliverUserCreatedResponse(w, response)
	return
}
