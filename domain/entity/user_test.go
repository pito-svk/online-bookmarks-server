package entity

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHexID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hexId := generateHexID()

		decimalStringId := make([]byte, hex.DecodedLen(len(hexId)))

		_, err := hex.Decode(decimalStringId, []byte(hexId))

		assert.NoError(t, err)
	})
}

func TestGenerateId(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		id := generateID()

		assert.NotEmpty(t, id)
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		password := "randomPassword"
		hashedPassword, err := hashPassword(password)

		assert.NoError(t, err)
		assert.NotEqual(t, password, hashedPassword)
	})
}

func TestSetID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := User{
			Email:     "random@example.com",
			Password:  "hashedPassword",
			FirstName: "John",
			LastName:  "Doe",
		}

		user.SetID()

		assert.NotEmpty(t, user.ID)
	})
}

func TestSetHashedPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := User{
			ID:        "5f555a4686dbe11fc3adbb9b",
			Email:     "random@example.com",
			Password:  "hashedPassword",
			FirstName: "John",
			LastName:  "Doe",
		}

		err := user.SetHashedPassword()

		assert.NoError(t, err)
		assert.NotEqual(t, user.Password, "hashedPassword")
	})
}

func TestClearPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := User{
			ID:        "5f555a4686dbe11fc3adbb9b",
			Email:     "random@example.com",
			Password:  "hashedPassword",
			FirstName: "John",
			LastName:  "Doe",
		}

		user.ClearPassword()

		assert.Empty(t, user.Password)
	})
}
