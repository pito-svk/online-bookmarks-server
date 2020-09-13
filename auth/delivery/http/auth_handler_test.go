package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/domain/mocks"
)

func TestLowercaseFirstLetter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assert.Equal(t, lowercaseFirstLetter("Email"), "email")
		assert.Equal(t, lowercaseFirstLetter("FirstName"), "firstName")
	})
}

func TestValidateCreateUserInput(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userData := userDataInput{
			Email:     "random@example.com",
			Password:  "demouser",
			FirstName: "John",
			LastName:  "Doe",
		}

		err := validateCreateUserInput(&userData)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		userData := userDataInput{
			Email:     "invalidEmail.com",
			Password:  "demouser",
			FirstName: "John",
			LastName:  "Doe",
		}

		err := validateCreateUserInput(&userData)

		assert.Error(t, err)
	})
}

func TestDeliverUserCreatedResponse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()

		user := userCreatedResponse{
			ID:        "5f555a4686dbe11fc3adbb9b",
			Email:     "random@example.com",
			FirstName: "John",
			LastName:  "Doe",
			AuthData: authData{
				Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjVmNWIyOTA5NjMwZTdmMmQ5NWU5MjZkMCJ9.CxxXHpVzS5f0Psl34iLXR9sg3HCEB0dYglMfhvWHoZ4",
			},
		}

		err := deliverUserCreatedResponse(w, user)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusCreated, w.Code)

		var httpResponse userCreatedResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&httpResponse); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, user, httpResponse)
	})
}

func TestRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUsecase := mocks.NewAuthUsecase(mockUserRepo)
	mockLogger := mocks.NewLogger()
	jwtSecret := "SECRET"

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()

		userData := userDataInput{
			Email:     "random@example.com",
			Password:  "demouser",
			FirstName: "John",
			LastName:  "Doe",
		}

		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			panic(err)
		}

		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(string(userDataJSON)))

		handler := AuthHandler{
			AuthUsecase: mockUsecase,
			Logger:      mockLogger,
			JwtSecret:   jwtSecret,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		err = json.Unmarshal(w.Body.Bytes(), &jsonResponse)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NotEmpty(t, jsonResponse["id"])
		assert.Empty(t, jsonResponse["password"])
		assert.Equal(t, "random@example.com", jsonResponse["email"])
		assert.Equal(t, "John", jsonResponse["firstName"])
		assert.Equal(t, "Doe", jsonResponse["lastName"])

		authToken := jsonResponse["authData"].(map[string]interface{})["token"].(string)

		assert.NotEmpty(t, authToken)

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		assert.NoError(t, err)

		claims, ok := token.Claims.(jwt.MapClaims)

		assert.True(t, ok)
		assert.NoError(t, claims.Valid())

		assert.Equal(t, jsonResponse["id"], claims["id"])
	})

	t.Run("duplicate", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{ "email": "random@example.com", "password": "demouser", "firstName": "John", "lastName": "Doe" }`))

		handler := AuthHandler{
			AuthUsecase: mockUsecase,
			Logger:      mockLogger,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Equal(t, "User already exists", jsonResponse["error"])
	})

	t.Run("missing email", func(t *testing.T) {
		w := httptest.NewRecorder()

		userData := userDataInput{
			Password:  "demouser",
			FirstName: "John",
			LastName:  "Doe",
		}

		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			panic(err)
		}

		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(string(userDataJSON)))

		handler := AuthHandler{
			AuthUsecase: mockUsecase,
			Logger:      mockLogger,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing email", jsonResponse["error"])
	})

	t.Run("missing password", func(t *testing.T) {
		w := httptest.NewRecorder()

		userData := userDataInput{
			Email:     "random@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}

		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			panic(err)
		}

		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(string(userDataJSON)))

		handler := AuthHandler{
			AuthUsecase: mockUsecase,
			Logger:      mockLogger,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing password", jsonResponse["error"])
	})

	t.Run("missing firstName", func(t *testing.T) {
		w := httptest.NewRecorder()

		userData := userDataInput{
			Email:    "random@example.com",
			Password: "demouser",
			LastName: "Doe",
		}

		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			panic(err)
		}

		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(string(userDataJSON)))

		handler := AuthHandler{
			AuthUsecase: mockUsecase,
			Logger:      mockLogger,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing firstName", jsonResponse["error"])
	})

	t.Run("missing lastName", func(t *testing.T) {
		w := httptest.NewRecorder()

		userData := userDataInput{
			Email:     "random@example.com",
			Password:  "demouser",
			FirstName: "John",
		}

		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			panic(err)
		}

		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(string(userDataJSON)))

		handler := AuthHandler{
			AuthUsecase: mockUsecase,
			Logger:      mockLogger,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing lastName", jsonResponse["error"])
	})

	t.Run("invalid email", func(t *testing.T) {
		w := httptest.NewRecorder()

		userData := userDataInput{
			Email:     "invalidEmail",
			Password:  "demouser",
			FirstName: "John",
			LastName:  "Doe",
		}

		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			panic(err)
		}

		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(string(userDataJSON)))

		handler := AuthHandler{
			AuthUsecase: mockUsecase,
			Logger:      mockLogger,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid email", jsonResponse["error"])
	})

	t.Run("invalid attribute type", func(t *testing.T) {
		w := httptest.NewRecorder()

		userDataJSONString := `{ "email": "random@example.com", "password": "demouser", "firstName": 1, "lastName": "Doe" }`

		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(userDataJSONString))

		handler := AuthHandler{
			AuthUsecase: mockUsecase,
			Logger:      mockLogger,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Error parsing JSON body", jsonResponse["error"])
	})
}
