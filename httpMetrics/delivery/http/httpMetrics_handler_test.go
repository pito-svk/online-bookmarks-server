package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/domain/mocks"
)

func TestLogHTTPMetrics(t *testing.T) {
	mockUsecase := mocks.NewHTTPMetricsUsecase()
	mockLogger := mocks.NewLogger()

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest("POST", "/users/register", strings.NewReader(""))
		r.Header.Set("Referer", "https://www.example.com")
		r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; SM-G935S Build/MMB29K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36")
		r.Header.Set("X-Forwarded-For", "192.168.2.1 , 217.73.23.164")

		handler := HTTPMetricsHandler{
			HTTPMetricsUsecase: mockUsecase,
			Logger:             mockLogger,
		}

		genericHTTPHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		})
		httpMetricsHandler := handler.LogHTTPMetrics(genericHTTPHandler)

		httpMetricsHandler.ServeHTTP(w, r)

		requestMetrics := handler.HTTPMetrics.RequestMetrics
		responseMetrics := handler.HTTPMetrics.ResponseMetrics

		assert.Equal(t, "/users/register", requestMetrics.URI)
		assert.Equal(t, "POST", requestMetrics.Method)
		assert.Equal(t, "https://www.example.com", requestMetrics.Referer)
		assert.Equal(t, "Mozilla/5.0 (Linux; Android 6.0.1; SM-G935S Build/MMB29K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36", requestMetrics.UserAgent)
		assert.Equal(t, "217.73.23.164", requestMetrics.IP)
		assert.Equal(t, http.StatusCreated, responseMetrics.Code)
	})

	t.Run("success 2", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/users/register", strings.NewReader(""))
		r.Header.Set("Referer", "https://www.example.com/example")
		r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.1.1; G8231 Build/41.2.A.0.219; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36")
		r.Header.Set("X-Real-Ip", "217.73.23.164")

		handler := HTTPMetricsHandler{
			HTTPMetricsUsecase: mockUsecase,
			Logger:             mockLogger,
		}

		genericHTTPHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})
		httpMetricsHandler := handler.LogHTTPMetrics(genericHTTPHandler)

		httpMetricsHandler.ServeHTTP(w, r)

		requestMetrics := handler.HTTPMetrics.RequestMetrics
		responseMetrics := handler.HTTPMetrics.ResponseMetrics

		assert.Equal(t, "/users/register", requestMetrics.URI)
		assert.Equal(t, "POST", requestMetrics.Method)
		assert.Equal(t, "https://www.example.com/example", requestMetrics.Referer)
		assert.Equal(t, "Mozilla/5.0 (Linux; Android 7.1.1; G8231 Build/41.2.A.0.219; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36", requestMetrics.UserAgent)
		assert.Equal(t, "217.73.23.164", requestMetrics.IP)
		assert.Equal(t, http.StatusBadRequest, responseMetrics.Code)
	})

	t.Run("success 3", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/users/register", strings.NewReader(""))
		r.Header.Set("Referer", "https://www.example.com/example")
		r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.1.1; G8231 Build/41.2.A.0.219; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36")
		r.RemoteAddr = "217.73.23.163"

		handler := HTTPMetricsHandler{
			HTTPMetricsUsecase: mockUsecase,
			Logger:             mockLogger,
		}

		genericHTTPHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		httpMetricsHandler := handler.LogHTTPMetrics(genericHTTPHandler)

		httpMetricsHandler.ServeHTTP(w, r)

		requestMetrics := handler.HTTPMetrics.RequestMetrics
		responseMetrics := handler.HTTPMetrics.ResponseMetrics

		assert.Equal(t, "/users/register", requestMetrics.URI)
		assert.Equal(t, "POST", requestMetrics.Method)
		assert.Equal(t, "https://www.example.com/example", requestMetrics.Referer)
		assert.Equal(t, "Mozilla/5.0 (Linux; Android 7.1.1; G8231 Build/41.2.A.0.219; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36", requestMetrics.UserAgent)
		assert.Equal(t, "217.73.23.163", requestMetrics.IP)
		assert.Equal(t, http.StatusInternalServerError, responseMetrics.Code)
	})
}
