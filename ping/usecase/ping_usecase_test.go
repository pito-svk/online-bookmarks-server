package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		u := NewPingUsecase()

		pingResponse := u.Get()

		assert.Equal(t, "PONG", pingResponse)
	})
}
