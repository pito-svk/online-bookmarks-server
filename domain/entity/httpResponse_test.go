package entity

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewResponseWriterWithMetrics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_w := httptest.NewRecorder()
		w := NewResponseWriterWithMetrics(_w)

		assert.Equal(t, http.StatusOK, w.StatusCode)
		assert.Equal(t, 0, w.Duration)
		assert.IsType(t, time.Time{}, w.requestTimeStart)
		assert.IsType(t, _w, w.ResponseWriter)

		w.WriteHeader(400)

		assert.Equal(t, 400, w.StatusCode)
	})
}

func TestCalcRequestDuration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		requestStart := time.Now()
		now := requestStart.Add(time.Millisecond * 50)

		duration := calcRequestDuration(requestStart, now)

		assert.Equal(t, 50, duration)
	})

	t.Run("success 2", func(t *testing.T) {
		requestStart := time.Now()

		duration := calcRequestDuration(requestStart, requestStart)

		assert.Equal(t, 1, duration)
	})
}

func TestServeHTTP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_w := httptest.NewRecorder()
		w := NewResponseWriterWithMetrics(_w)

		r := httptest.NewRequest("GET", "/ping", strings.NewReader(""))

		handler := HTTPHandlerSettingRequestDuration{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		}

		handler.ServeHTTP(w, r)

		assert.Equal(t, 1, w.Duration)
	})

	t.Run("success 2", func(t *testing.T) {
		_w := httptest.NewRecorder()
		w := NewResponseWriterWithMetrics(_w)

		w.requestTimeStart = time.Now().Add(-time.Millisecond * 5)

		r := httptest.NewRequest("GET", "/ping", strings.NewReader(""))

		handler := HTTPHandlerSettingRequestDuration{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		}

		handler.ServeHTTP(w, r)

		assert.Equal(t, 5, w.Duration)
	})
}

func TestDeliverHTTPErrors(t *testing.T) {
	t.Run("error parsing JSON body error", func(t *testing.T) {
		w := httptest.NewRecorder()

		DeliverErrorParsingJSONBodyHTTPError(w)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request error", func(t *testing.T) {
		w := httptest.NewRecorder()

		err := errors.New("Bad request")

		DeliverBadRequestHTTPError(w, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}