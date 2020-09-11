package entity

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPrivateIpAddress(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ip := IPAddress{Address: "127.0.0.1"}

		isPrivateIP, err := ip.IsPrivate()

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 2", func(t *testing.T) {
		ip := IPAddress{Address: "192.168.0.0"}

		isPrivateIP, err := ip.IsPrivate()

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 3", func(t *testing.T) {
		ip := IPAddress{Address: "192.168.255.255"}

		isPrivateIP, err := ip.IsPrivate()

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 4", func(t *testing.T) {
		ip := IPAddress{Address: "172.16.0.0"}

		isPrivateIP, err := ip.IsPrivate()

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 5", func(t *testing.T) {
		ip := IPAddress{Address: "172.31.255.255"}

		isPrivateIP, err := ip.IsPrivate()

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 6", func(t *testing.T) {
		ip := IPAddress{Address: "10.0.0.0"}

		isPrivateIP, err := ip.IsPrivate()

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 7", func(t *testing.T) {
		ip := IPAddress{Address: "10.255.255.255"}

		isPrivateIP, err := ip.IsPrivate()

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 9 (public ip address)", func(t *testing.T) {
		ip := IPAddress{Address: "217.73.23.164"}

		isPrivateIP, err := ip.IsPrivate()

		assert.NoError(t, err)
		assert.Equal(t, false, isPrivateIP)
	})

	t.Run("success 10", func(t *testing.T) {
		ip := IPAddress{Address: "127.0.0.0"}

		isPrivateIP, err := ip.IsPrivate()

		assert.NoError(t, err)
		assert.Equal(t, true, isPrivateIP)
	})

	t.Run("success 11", func(t *testing.T) {
		ip := IPAddress{Address: "169.254.0.0"}

		isPrivateIP, err := ip.IsPrivate()

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

func TestParseIPFromXRealIPHeader(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ipAddress := parseIPFromXRealIPHeader("217.73.23.164")

		assert.Equal(t, "217.73.23.164", ipAddress)
	})
}

func TestParseRemoteAddr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assert.Equal(t, "217.73.23.164", parseRemoteIPAddress("217.73.23.164"))
		assert.Equal(t, "217.73.23.164", parseRemoteIPAddress("217.73.23.164:3000"))
	})
}

func TestGetIpAddressFromHTTPRequest(t *testing.T) {
	t.Run("success with x-forwarded-for", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.Header.Set("X-Forwarded-For", "192.168.2.1 , 217.73.23.164")

		ipAddress, err := GetIPAddressFromHTTPRequest(r)

		assert.NoError(t, err)
		assert.Equal(t, "217.73.23.164", ipAddress)
	})

	t.Run("success with x-real-ip", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.Header.Set("X-Real-Ip", "217.73.23.164")

		ipAddress, err := GetIPAddressFromHTTPRequest(r)

		assert.NoError(t, err)
		assert.Equal(t, "217.73.23.164", ipAddress)
	})

	t.Run("success with remoteAddr", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.RemoteAddr = "217.73.23.164"

		ipAddress, err := GetIPAddressFromHTTPRequest(r)

		assert.NoError(t, err)
		assert.Equal(t, "217.73.23.164", ipAddress)
	})

	t.Run("success with remoteAddr with port", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(""))
		r.RemoteAddr = "217.73.23.164:3000"

		ipAddress, err := GetIPAddressFromHTTPRequest(r)

		assert.NoError(t, err)
		assert.Equal(t, "217.73.23.164", ipAddress)
	})
}
