package usecase_test

import (
	"encoding/hex"
	"testing"

	"github.com/dgrijalva/jwt-go"
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

func TestGenerateId(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		id := usecase.GenerateID()

		assert.NotEmpty(t, id)
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		password := "randomPassword"
		hashedPassword, err := usecase.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEqual(t, password, hashedPassword)
	})
}

func TestComparePasswords(t *testing.T) {
	t.Run("password match", func(t *testing.T) {
		hashedPassword := "$2a$10$F6U1i4rENwIWiZXBigkxTe3bXV.WixvL10MSWRY6e9icQXuXWmT5."
		password := "randomPassword"

		passwordMatch := usecase.ComparePasswords(hashedPassword, password)

		assert.True(t, passwordMatch)
	})

	t.Run("password match 2", func(t *testing.T) {
		hashedPassword := "$2a$10$cbHkZbUT513CN4aIKVAq8OEs1QBr/qk562NiNtVQ.rakI/qrzxSBi"
		password := "randomPassword"

		passwordMatch := usecase.ComparePasswords(hashedPassword, password)

		assert.True(t, passwordMatch)
	})

	t.Run("no match", func(t *testing.T) {
		hashedPassword := "$2a$10$F6U1i4rENwIWiZXBigkxTe3bXV.WixvL10MSWRY6e9icQXuXWmT5."
		password := "wrongPassword"

		passwordMatch := usecase.ComparePasswords(hashedPassword, password)

		assert.False(t, passwordMatch)
	})
}

func TestGenerateAuthToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		jwtSecret := "SECRET"

		claimData := map[string]interface{}{
			"id": "5f5410bd3cfca9b341bdfe4c",
		}

		authToken, err := usecase.GenerateAuthToken(claimData, jwtSecret)

		assert.NoError(t, err)
		assert.NotEmpty(t, authToken)

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		assert.NoError(t, err)

		claims, ok := token.Claims.(jwt.MapClaims)

		assert.True(t, ok)
		assert.NoError(t, claims.Valid())

		assert.Equal(t, claims["id"], "5f5410bd3cfca9b341bdfe4c")
	})
}

func TestRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		u := usecase.NewAuthUsecase(mockUserRepo)

		userData := entity.User{
			Email:     "random@example.com",
			Password:  "securePassword",
			FirstName: "John",
			LastName:  "Doe",
		}

		duplicateUserData := userData

		registeredUser, err := u.Register(&userData)

		assert.NoError(t, err)

		assert.NotEmpty(t, registeredUser.ID)
		assert.Equal(t, "random@example.com", registeredUser.Email)
		assert.Empty(t, registeredUser.Password)
		assert.Equal(t, "John", registeredUser.FirstName)
		assert.Equal(t, "Doe", registeredUser.LastName)

		secondUserData := entity.User{
			Email:     "random2@example.com",
			Password:  "securePassword",
			FirstName: "Martin",
			LastName:  "Appleseed",
		}

		secondRegisteredUser, err := u.Register(&secondUserData)

		assert.NoError(t, err)

		assert.NotEqual(t, registeredUser.ID, secondRegisteredUser.ID)
		assert.Equal(t, "random2@example.com", secondRegisteredUser.Email)
		assert.Empty(t, secondRegisteredUser.Password)
		assert.Equal(t, "Martin", secondRegisteredUser.FirstName)
		assert.Equal(t, "Appleseed", secondRegisteredUser.LastName)

		duplicateUser, err := u.Register(&duplicateUserData)

		assert.Error(t, err)
		assert.Nil(t, duplicateUser)
	})
}
