package mocks

import (
	"peterparada.com/online-bookmarks/auth/usecase"
	"peterparada.com/online-bookmarks/domain/entity"
)

type AuthUsecase struct {
}

func (a *AuthUsecase) Register(user *entity.User) (*entity.User, error) {
	user.ID = usecase.GenerateID()
	hashedPassword, _ := usecase.HashPassword(user.Password)

	user.Password = hashedPassword

	return user, nil
}
