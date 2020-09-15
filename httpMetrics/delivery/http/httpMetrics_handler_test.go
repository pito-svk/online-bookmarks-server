package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/domain/entity"
	"peterparada.com/online-bookmarks/domain/mocks"
)

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"peterparada.com/online-bookmarks/domain/entity"
// 	"peterparada.com/online-bookmarks/domain/mocks"
// )

// func TestGetHTTPRequestMetrics(t *testing.T) {
// 	t.Run("success", func(t *testing.T) {
// 		u := NewHTTPMetricsUsecase()

// 		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
// 		r.Header.Set("Referer", "https://www.example.com")
// 		r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; SM-G935S Build/MMB29K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36")
// 		r.Header.Set("X-Forwarded-For", "192.168.2.1 , 217.73.23.164")

// 		requestMetrics, err := u.GetHTTPRequestMetrics(r)

// 		assert.NoError(t, err)
// 		assert.Equal(t, "/auth/register", requestMetrics.URI)
// 		assert.Equal(t, "POST", requestMetrics.Method)
// 		assert.Equal(t, "https://www.example.com", requestMetrics.Referer)
// 		assert.Equal(t, "Mozilla/5.0 (Linux; Android 6.0.1; SM-G935S Build/MMB29K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36", requestMetrics.UserAgent)
// 		assert.Equal(t, "217.73.23.164", requestMetrics.IP)
// 	})

// 	t.Run("success 2", func(t *testing.T) {
// 		u := NewHTTPMetricsUsecase()

// 		r := httptest.NewRequest("GET", "/ping", strings.NewReader(""))
// 		r.Header.Set("X-Real-Ip", "217.73.23.164")

// 		requestMetrics, err := u.GetHTTPRequestMetrics(r)

// 		assert.NoError(t, err)
// 		assert.Equal(t, "/ping", requestMetrics.URI)
// 		assert.Equal(t, "GET", requestMetrics.Method)
// 		assert.Equal(t, "", requestMetrics.Referer)
// 		assert.Equal(t, "", requestMetrics.UserAgent)
// 		assert.Equal(t, "217.73.23.164", requestMetrics.IP)
// 	})

// 	t.Run("success 3", func(t *testing.T) {
// 		u := NewHTTPMetricsUsecase()

// 		r := httptest.NewRequest("GET", "/users?queryString=123", strings.NewReader(""))
// 		r.Header.Set("Referer", "https://www.example.com")
// 		r.RemoteAddr = "217.73.23.163"

// 		requestMetrics, err := u.GetHTTPRequestMetrics(r)

// 		assert.NoError(t, err)
// 		assert.Equal(t, "/users?queryString=123", requestMetrics.URI)
// 		assert.Equal(t, "GET", requestMetrics.Method)
// 		assert.Equal(t, "https://www.example.com", requestMetrics.Referer)
// 		assert.Equal(t, "217.73.23.163", requestMetrics.IP)
// 	})
// }

// func TestGetHTTPResponseMetrics(t *testing.T) {
// 	t.Run("success", func(t *testing.T) {
// 		u := NewHTTPMetricsUsecase()

// 		_w := httptest.NewRecorder()

// 		r := httptest.NewRequest("POST", "/users/register", strings.NewReader(""))
// 		w := entity.NewResponseWriterWithMetrics(_w)

// 		httpHandler := mocks.HTTPHandlerSettingRequestDuration{
// 			Handler:           http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
// 			WriterWithMetrics: w,
// 			Duration:          123,
// 		}

// 		httpHandler.ServeHTTP(w, r)

// 		responseMetrics := u.GetHTTPResponseMetrics(w)

// 		assert.Equal(t, http.StatusOK, responseMetrics.Code)
// 		assert.Equal(t, 123, responseMetrics.Duration)
// 	})

// 	t.Run("success 2", func(t *testing.T) {
// 		u := NewHTTPMetricsUsecase()

// 		_w := httptest.NewRecorder()

// 		r := httptest.NewRequest("POST", "/users/register", strings.NewReader(""))
// 		w := entity.NewResponseWriterWithMetrics(_w)
// 		w.WriteHeader(http.StatusNotFound)

// 		httpHandler := mocks.HTTPHandlerSettingRequestDuration{
// 			Handler:           http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
// 			WriterWithMetrics: w,
// 			Duration:          60,
// 		}

// 		httpHandler.ServeHTTP(w, r)

// 		responseMetrics := u.GetHTTPResponseMetrics(w)

// 		assert.Equal(t, http.StatusNotFound, responseMetrics.Code)
// 		assert.Equal(t, 60, responseMetrics.Duration)
// 	})
// }

func TestLogHTTPMetrics(t *testing.T) {
	mockUsecase := mocks.NewHTTPMetricsUsecase()
	mockLogger := mocks.NewLogger()

	t.Run("success", func(t *testing.T) {
		_w := httptest.NewRecorder()
		w := entity.NewResponseWriterWithMetrics(_w)

		r := httptest.NewRequest("POST", "/users/register", strings.NewReader(""))
		r.Header.Set("Referer", "https://www.example.com")
		r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; SM-G935S Build/MMB29K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36")
		r.Header.Set("X-Forwarded-For", "192.168.2.1 , 217.73.23.164")

		httpHandler := &mocks.HTTPHandlerSettingRequestDuration{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			}),
			WriterWithMetrics: w,
			Duration:          136,
		}

		handler := HTTPMetricsHandler{
			HTTPMetricsUsecase: mockUsecase,
			Logger:             mockLogger,
			HTTPHandler:        httpHandler,
		}

		httpMetricsHandler := handler.LogHTTPMetrics(httpHandler.Handler)

		httpMetricsHandler.ServeHTTP(w, r)

		requestMetrics := handler.HTTPMetrics.RequestMetrics
		responseMetrics := handler.HTTPMetrics.ResponseMetrics

		assert.Equal(t, "/users/register", requestMetrics.URI)
		assert.Equal(t, "POST", requestMetrics.Method)
		assert.Equal(t, "https://www.example.com", requestMetrics.Referer)
		assert.Equal(t, "Mozilla/5.0 (Linux; Android 6.0.1; SM-G935S Build/MMB29K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36", requestMetrics.UserAgent)
		assert.Equal(t, "217.73.23.164", requestMetrics.IP)
		assert.Equal(t, http.StatusCreated, responseMetrics.Code)
		assert.Equal(t, 136, responseMetrics.Duration)
	})

	t.Run("success 2", func(t *testing.T) {
		_w := httptest.NewRecorder()
		w := entity.NewResponseWriterWithMetrics(_w)

		r := httptest.NewRequest("POST", "/users/register", strings.NewReader(""))
		r.Header.Set("Referer", "https://www.example.com/example")
		r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.1.1; G8231 Build/41.2.A.0.219; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36")
		r.Header.Set("X-Real-Ip", "217.73.23.164")

		httpHandler := &mocks.HTTPHandlerSettingRequestDuration{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}),
			WriterWithMetrics: w,
			Duration:          55,
		}

		handler := HTTPMetricsHandler{
			HTTPMetricsUsecase: mockUsecase,
			Logger:             mockLogger,
			HTTPHandler:        httpHandler,
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
		assert.Equal(t, 55, responseMetrics.Duration)
	})

	t.Run("success 3", func(t *testing.T) {
		_w := httptest.NewRecorder()
		w := entity.NewResponseWriterWithMetrics(_w)

		r := httptest.NewRequest("POST", "/users/register", strings.NewReader(""))
		r.Header.Set("Referer", "https://www.example.com/example")
		r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.1.1; G8231 Build/41.2.A.0.219; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36")
		r.RemoteAddr = "217.73.23.163"

		httpHandler := &mocks.HTTPHandlerSettingRequestDuration{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			WriterWithMetrics: w,
			Duration:          62,
		}

		handler := HTTPMetricsHandler{
			HTTPMetricsUsecase: mockUsecase,
			Logger:             mockLogger,
			HTTPHandler:        httpHandler,
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
		assert.Equal(t, 62, responseMetrics.Duration)
	})
}
