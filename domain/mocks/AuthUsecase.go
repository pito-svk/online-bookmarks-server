package mocks

import (
	"peterparada.com/online-bookmarks/auth/usecase"
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type AuthUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthUsecase(userRepo domain.UserRepository) domain.AuthUsecase {
	return &AuthUsecase{
		userRepo,
	}
}

func (a *AuthUsecase) Register(user *entity.User) (*entity.User, error) {
	user.ID = usecase.GenerateID()
	hashedPassword, _ := usecase.HashPassword(user.Password)

	user.Password = hashedPassword

	return a.userRepo.Store(user)
}
