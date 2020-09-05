package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	_authHttpDelivery "peterparada.com/online-bookmarks/auth/delivery/http"
	"peterparada.com/online-bookmarks/domain/mocks"
)

func TestRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUsecase := mocks.NewAuthUsecase(mockUserRepo)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{ "email": "random@example.com", "password": "demouser", "firstName": "John", "lastName": "Doe" }`))

		handler := _authHttpDelivery.AuthHandler{
			AuthUsecase: mockUsecase,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NotEmpty(t, jsonResponse["id"])
		assert.Empty(t, jsonResponse["password"])
		assert.Equal(t, "random@example.com", jsonResponse["email"])
		assert.Equal(t, "John", jsonResponse["firstName"])
		assert.Equal(t, "Doe", jsonResponse["lastName"])
	})

	t.Run("duplicate", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{ "email": "random@example.com", "password": "demouser", "firstName": "John", "lastName": "Doe" }`))

		handler := _authHttpDelivery.AuthHandler{
			AuthUsecase: mockUsecase,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Equal(t, "Email already exists", jsonResponse["error"])
	})

	t.Run("missing email", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{ "password": "demouser", "firstName": "John", "lastName": "Doe" }`))

		handler := _authHttpDelivery.AuthHandler{
			AuthUsecase: mockUsecase,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing email", jsonResponse["error"])
	})

	t.Run("missing password", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{ "email": "random@example.com", "firstName": "John", "lastName": "Doe" }`))

		handler := _authHttpDelivery.AuthHandler{
			AuthUsecase: mockUsecase,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing password", jsonResponse["error"])
	})

	t.Run("missing firstName", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{ "email": "random@example.com", "password": "demouser", "lastName": "Doe" }`))

		handler := _authHttpDelivery.AuthHandler{
			AuthUsecase: mockUsecase,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing firstName", jsonResponse["error"])
	})

	t.Run("missing lastName", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{ "email": "random@example.com", "password": "demouser", "firstName": "John" }`))

		handler := _authHttpDelivery.AuthHandler{
			AuthUsecase: mockUsecase,
		}

		handler.RegisterUser(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing lastName", jsonResponse["error"])
	})
}
