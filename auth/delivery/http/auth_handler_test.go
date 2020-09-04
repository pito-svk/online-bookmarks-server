package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	_authHttpDelivery "peterparada.com/online-bookmarks/auth/delivery/http"
	"peterparada.com/online-bookmarks/domain/mocks"
)

func TestRegister(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{ "email": "random@example.com", "password": "demouser", "firstName": "John", "lastName": "Doe" }`))

	mockUsecase := new(mocks.AuthUsecase)

	handler := _authHttpDelivery.AuthHandler{
		AuthUsecase: mockUsecase,
	}

	handler.RegisterUser(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
}
