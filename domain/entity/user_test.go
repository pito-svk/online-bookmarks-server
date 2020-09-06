package entity_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/domain/entity"
)

func TestGenerateHexID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hexId := entity.GenerateHexID()

		decimalStringId := make([]byte, hex.DecodedLen(len(hexId)))

		_, err := hex.Decode(decimalStringId, []byte(hexId))

		assert.NoError(t, err)
	})
}

func TestGenerateId(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		id := entity.GenerateID()

		assert.NotEmpty(t, id)
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		password := "randomPassword"
		hashedPassword, err := entity.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEqual(t, password, hashedPassword)
	})
}
