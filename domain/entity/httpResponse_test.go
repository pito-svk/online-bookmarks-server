package entity

import (
	"encoding/json"
	"errors"
	"log"
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

		var httpResponse httpErrorMessage
		if err := json.NewDecoder(w.Result().Body).Decode(&httpResponse); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, "Error parsing JSON body", httpResponse.Error)
	})

	t.Run("bad request error", func(t *testing.T) {
		w := httptest.NewRecorder()

		err := errors.New("Bad request")

		DeliverBadRequestHTTPError(w, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var httpResponse httpErrorMessage
		if err := json.NewDecoder(w.Result().Body).Decode(&httpResponse); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, "Bad request", httpResponse.Error)
	})

	t.Run("conflict error", func(t *testing.T) {
		w := httptest.NewRecorder()

		err := errors.New("User already exists")

		DeliverConflictHTTPError(w, err)

		assert.Equal(t, http.StatusConflict, w.Code)

		var httpResponse httpErrorMessage
		if err := json.NewDecoder(w.Result().Body).Decode(&httpResponse); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, "User already exists", httpResponse.Error)
	})

	t.Run("internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()

		DeliverInternalServerErrorHTTPError(w)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var httpResponse httpErrorMessage
		if err := json.NewDecoder(w.Result().Body).Decode(&httpResponse); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, "Internal Server Error", httpResponse.Error)
	})

	t.Run("endpoint not found", func(t *testing.T) {
		w := httptest.NewRecorder()

		DeliverEndpointNotFoundError(w)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var httpResponse httpErrorMessage
		if err := json.NewDecoder(w.Result().Body).Decode(&httpResponse); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, "Endpoint not found", httpResponse.Error)
	})
}
