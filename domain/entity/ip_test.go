package entity

import (
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
