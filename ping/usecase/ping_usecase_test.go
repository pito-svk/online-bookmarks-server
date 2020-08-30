package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/ping/usecase"
)

func TestGet(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		u := usecase.NewPingUsecase()

		pingResponse := u.Get()

		assert.Equal(t, pingResponse, "PONG")
	})
}
