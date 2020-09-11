package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"peterparada.com/online-bookmarks/domain/mocks"
)

func TestLogHTTPMetrics(t *testing.T) {
	mockUsecase := mocks.NewHTTPMetricsUsecase()
	mockLogger := mocks.NewLogger()

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/ping", strings.NewReader(""))

		handler := HTTPMetricsHandler{
			HTTPMetricsUsecase: mockUsecase,
			Logger:             mockLogger,
		}

		genericHTTPHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		httpMetricsHandler := handler.LogHTTPMetrics(genericHTTPHandler)

		httpMetricsHandler.ServeHTTP(w, r)
	})
}
