package http

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/stretchr/testify/assert"
)

func TestIsPrivateIpAddress(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("127.0.0.1")

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 2", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("192.168.0.0")

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 3", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("192.168.255.255")

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 4", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("172.16.0.0")

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 5", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("172.31.255.255")

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 6", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("10.0.0.0")

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 7", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("10.255.255.255")

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 9 (public ip address)", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("217.73.23.164")

		assert.NoError(t, err)
		assert.Equal(t, false, isPrivateIP)
	})

	t.Run("success 10", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("127.0.0.0")

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 11", func(t *testing.T) {
		isPrivateIP, err := isPrivateIPAddress("169.254.0.0")

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})
}

func TestParseIPFromXForwardedForHeader(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ipAddress, err := parseIPFromXForwardedForHeader("192.168.2.1 , 217.73.23.164")

		assert.NoError(t, err)
		assert.Equal(t, "217.73.23.164", ipAddress)
	})

	t.Run("success 2", func(t *testing.T) {
		ipAddress, err := parseIPFromXForwardedForHeader("217.73.23.165")

		assert.NoError(t, err)
		assert.Equal(t, "217.73.23.165", ipAddress)
	})

	t.Run("empty", func(t *testing.T) {
		ipAddress, err := parseIPFromXForwardedForHeader("192.168.2.1")

		assert.NoError(t, err)
		assert.Equal(t, "", ipAddress)
	})
}

func TestParseRemoteAddr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assert.Equal(t, "217.73.23.164", parseRemoteIPAddress("217.73.23.164"))
		assert.Equal(t, "217.73.23.164", parseRemoteIPAddress("217.73.23.164:3000"))
	})
}

func TestGetIpAddressFromHttpRequest(t *testing.T) {
	t.Run("success with x-forwarded-for", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.Header.Set("X-Forwarded-For", "192.168.2.1 , 217.73.23.164")

		ipAddress, err := getIPAddressFromHttpRequest(r)

		assert.NoError(t, err)
		assert.Equal(t, "217.73.23.164", ipAddress)
	})

	t.Run("success with x-real-ip", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.Header.Set("X-Real-Ip", "217.73.23.164")

		ipAddress, err := getIPAddressFromHttpRequest(r)

		assert.NoError(t, err)
		assert.Equal(t, "217.73.23.164", ipAddress)
	})

	t.Run("success with remoteAddr", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.RemoteAddr = "217.73.23.164"

		ipAddress, err := getIPAddressFromHttpRequest(r)

		assert.NoError(t, err)
		assert.Equal(t, "217.73.23.164", ipAddress)
	})

	t.Run("success with remoteAddr with port", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.RemoteAddr = "217.73.23.164:3000"

		ipAddress, err := getIPAddressFromHttpRequest(r)

		assert.NoError(t, err)
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

		httpRequestData, err := getHttpRequestData(r, httpResponseMetrics)

		assert.NoError(t, err)
		assert.Equal(t, "/auth/register", httpRequestData.URI)
		assert.Equal(t, "POST", httpRequestData.HTTPMethod)
		assert.Equal(t, "https://www.example.com", httpRequestData.Referer)
		assert.Equal(t, "Mozilla/5.0 (Linux; Android 6.0.1; SM-G935S Build/MMB29K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36", httpRequestData.UserAgent)
		assert.Equal(t, "217.73.23.164", httpRequestData.IP)
		assert.Equal(t, 400, httpRequestData.ResponseCode)
		assert.Equal(t, 233, httpRequestData.RequestDuration)
	})
}
