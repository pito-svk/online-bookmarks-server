package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleEndpointNotFound(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/notFound", strings.NewReader(""))

		handler := EndpointNotFoundHandler{}
		handler.HandleEndpointNotFound(w, r)

		var jsonResponse map[string]interface{}

		json.Unmarshal(w.Body.Bytes(), &jsonResponse)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "Endpoint not found", jsonResponse["error"])
	})
}
