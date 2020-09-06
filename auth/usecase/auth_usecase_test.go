package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/domain/entity"
	"peterparada.com/online-bookmarks/domain/mocks"
)

func TestRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		u := NewAuthUsecase(mockUserRepo)

		userData := entity.User{
			Email:     "random@example.com",
			Password:  "securePassword",
			FirstName: "John",
			LastName:  "Doe",
		}

		duplicateUserData := userData

		registeredUser, err := u.RegisterUser(&userData)

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

		secondRegisteredUser, err := u.RegisterUser(&secondUserData)

		assert.NoError(t, err)

		assert.NotEqual(t, registeredUser.ID, secondRegisteredUser.ID)
		assert.Equal(t, "random2@example.com", secondRegisteredUser.Email)
		assert.Empty(t, secondRegisteredUser.Password)
		assert.Equal(t, "Martin", secondRegisteredUser.FirstName)
		assert.Equal(t, "Appleseed", secondRegisteredUser.LastName)

		duplicateUser, err := u.RegisterUser(&duplicateUserData)

		assert.Error(t, err)
		assert.Nil(t, duplicateUser)
	})
}
