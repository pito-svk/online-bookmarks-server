package usecase_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/auth/usecase"
	"peterparada.com/online-bookmarks/domain/entity"
	"peterparada.com/online-bookmarks/domain/mocks"
)

func TestGenerateHexID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hexId := usecase.GenerateHexID()

		decimalStringId := make([]byte, hex.DecodedLen(len(hexId)))

		_, err := hex.Decode(decimalStringId, []byte(hexId))

		assert.NoError(t, err)
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		password := "randomPassword"
		hashedPassword := usecase.HashPassword(password)

		assert.NotEqual(t, password, hashedPassword)
	})
}

func TestRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		u := usecase.NewAuthUsecase(mockUserRepo)

		userData := entity.User{Email: "random@example.com", Password: "securePassword", FirstName: "John", LastName: "Doe"}
		duplicateUserData := userData

		registeredUser, err := u.Register(&userData)

		assert.NoError(t, err)

		assert.NotEmpty(t, registeredUser.ID)
		assert.Equal(t, registeredUser.Email, "random@example.com")
		assert.Empty(t, registeredUser.Password)
		assert.Equal(t, registeredUser.FirstName, "John")
		assert.Equal(t, registeredUser.LastName, "Doe")

		secondUserData := entity.User{Email: "random2@example.com", Password: "securePassword", FirstName: "Martin", LastName: "Appleseed"}

		secondRegisteredUser, err := u.Register(&secondUserData)

		assert.NoError(t, err)

		assert.NotEqual(t, registeredUser.ID, secondRegisteredUser.ID)
		assert.Equal(t, secondRegisteredUser.Email, "random2@example.com")
		assert.Empty(t, secondRegisteredUser.Password)
		assert.Equal(t, secondRegisteredUser.FirstName, "Martin")
		assert.Equal(t, secondRegisteredUser.LastName, "Appleseed")

		duplicateUser, err := u.Register(&duplicateUserData)

		assert.Error(t, err)
		assert.Nil(t, duplicateUser)
	})
}
