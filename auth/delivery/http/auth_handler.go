package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
	Logger      domain.Logger
	JwtSecret   string
}

func NewAuthHandler(router *chi.Mux, usecase domain.AuthUsecase, logger domain.Logger, jwtSecret string) {
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

func lowercaseFirstLetter(s string) string {
	var str strings.Builder

	firstLetter := string(s[0])
	restLetters := string(s[1:])

	str.WriteString(strings.ToLower(firstLetter))
	str.WriteString(restLetters)

	return str.String()
}

func validateCreateUserInput(userData *userDataInput) error {
	v := validator.New()

	err := v.Struct(userData)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		for _, err := range validationErrors {
			if err.Tag() == "required" {
				return fmt.Errorf("Missing %s", lowercaseFirstLetter(err.Field()))
			}

			if err.Tag() == "email" {
				return errors.New("Invalid email")
			}
		}
	}

	return err
}

func deliverUserCreatedResponse(w http.ResponseWriter, response userCreatedResponse) error {
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}

	return nil
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

func logErrorParsingJSONError(logger domain.Logger, err error) {
	logger.Error(map[string]interface{}{
		"method":  "auth_handler/RegisterUser",
		"pointer": "error_parsing_json",
	}, err)
}

func logInternalServerError(logger domain.Logger, err error, pointer string) {
	logger.Error(map[string]interface{}{
		"method":  "auth_handler/RegisterUser",
		"pointer": pointer,
	}, err)
}

func (authH *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	setJSONContentTypeInResponse(w)

	userData, err := parseUserDataFromRequestBody(r)
	if err != nil {
		logErrorParsingJSONError(authH.Logger, err)
		entity.DeliverErrorParsingJSONBodyHTTPError(w)
		return
	}
	err = validateCreateUserInput(userData)
	if err != nil {
		entity.DeliverBadRequestHTTPError(w, err)
		return
	}

	userObject := composeUserObjectFromUserData(userData)
	userResponse, err := authH.AuthUsecase.RegisterUser(&userObject)
	if err != nil {
		if err.Error() == "User already exists" {
			entity.DeliverConflictHTTPError(w, err)
		} else {
			logInternalServerError(authH.Logger, err, "register_user")
			entity.DeliverInternalServerErrorHTTPError(w)
		}
		return
	}
	authToken, err := authH.AuthUsecase.GenerateAuthToken(userResponse.ID, authH.JwtSecret)
	if err != nil {
		logInternalServerError(authH.Logger, err, "generate_auth_token")
		entity.DeliverInternalServerErrorHTTPError(w)
	}

	response := composeUserCreatedResponse(userResponse, authToken)

	err = deliverUserCreatedResponse(w, response)
	if err != nil {
		logInternalServerError(authH.Logger, err, "deliver_user_response")
		entity.DeliverInternalServerErrorHTTPError(w)
	}
}
