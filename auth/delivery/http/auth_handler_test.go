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

		assert.Equal(t, w.Code, http.StatusCreated)
		assert.NotEmpty(t, jsonResponse["id"])
		assert.Empty(t, jsonResponse["password"])
		assert.Equal(t, jsonResponse["email"], "random@example.com")
		assert.Equal(t, jsonResponse["firstName"], "John")
		assert.Equal(t, jsonResponse["lastName"], "Doe")
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

		assert.Equal(t, w.Code, http.StatusConflict)
		assert.Equal(t, jsonResponse["error"], "Email already exists")
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

		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Equal(t, jsonResponse["error"], "Missing email")
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

		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Equal(t, jsonResponse["error"], "Missing password")
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

		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Equal(t, jsonResponse["error"], "Missing firstName")
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

		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Equal(t, jsonResponse["error"], "Missing lastName")
	})
}
