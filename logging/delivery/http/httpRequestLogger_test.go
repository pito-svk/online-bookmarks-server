package http

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/stretchr/testify/assert"
)

func TestGetIpAddress(t *testing.T) {
	t.Run("success with x-forwarded-for", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.Header.Set("X-Forwarded-For", "217.73.23.164")

		ipAddress := getIpAddress(r)

		assert.Equal(t, "217.73.23.164", ipAddress)
	})

	t.Run("success with x-real-ip", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.Header.Set("X-Real-Ip", "217.73.23.164")

		ipAddress := getIpAddress(r)

		assert.Equal(t, "217.73.23.164", ipAddress)
	})
}

func TestGetHttpRequestData(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.Header.Set("Referer", "https://www.example.com")
		r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; SM-G935S Build/MMB29K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36")
		r.Header.Set("X-Forwarded-For", "217.73.23.164")

		httpResponseMetrics := httpsnoop.Metrics{
			Code:     400,
			Duration: time.Duration(233 * time.Millisecond),
		}

		httpRequestData := getHttpRequestData(r, httpResponseMetrics)

		assert.Equal(t, "/auth/register", httpRequestData.URI)
		assert.Equal(t, "POST", httpRequestData.HTTPMethod)
		assert.Equal(t, "https://www.example.com", httpRequestData.Referer)
		assert.Equal(t, "Mozilla/5.0 (Linux; Android 6.0.1; SM-G935S Build/MMB29K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36", httpRequestData.UserAgent)
		assert.Equal(t, "217.73.23.164", httpRequestData.IP)
		assert.Equal(t, 400, httpRequestData.ResponseCode)
		assert.Equal(t, 233, httpRequestData.RequestDuration)
	})
}
