package usecase

import (
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type authUsecase struct {
}

func NewAuthUsecase() domain.AuthUsecase {
	return &authUsecase{}
}

func (a *authUsecase) Register(u *entity.User) (error, entity.User) {
	return nil, entity.User{ID: "123", Email: "random@example.com", FirstName: "John", LastName: "Doe"}
}
