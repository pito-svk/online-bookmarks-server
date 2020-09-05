package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/domain/mocks"
	_pingHttpDelivery "peterparada.com/online-bookmarks/ping/delivery/http"
)

func TestGet(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ping", strings.NewReader(""))

	mockUsecase := new(mocks.PingUsecase)

	handler := _pingHttpDelivery.PingHandler{
		PingUsecase: mockUsecase,
	}

	handler.GetPing(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "PONG", w.Body.String())
}
