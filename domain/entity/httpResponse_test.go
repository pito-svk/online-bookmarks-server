package entity

import (
	"net/http"
	"net/http/httptest"
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