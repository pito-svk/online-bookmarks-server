package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/auth/usecase"
	"peterparada.com/online-bookmarks/domain/entity"
	"peterparada.com/online-bookmarks/domain/mocks"
)

func TestRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		u := usecase.NewAuthUsecase(mockUserRepo)

		userData := entity.User{Email: "random@example.com", Password: "securePassword", FirstName: "John", LastName: "Doe"}

		err, registeredUser := u.Register(&userData)

		assert.NoError(t, err)

		assert.NotEmpty(t, registeredUser.ID)
		assert.Equal(t, registeredUser.Email, "random@example.com")
		assert.Empty(t, registeredUser.Password)
		assert.Equal(t, registeredUser.FirstName, "John")
		assert.Equal(t, registeredUser.LastName, "Doe")
	})
}
