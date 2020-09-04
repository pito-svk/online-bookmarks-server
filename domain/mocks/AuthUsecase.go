package mocks

import "peterparada.com/online-bookmarks/domain/entity"

type AuthUsecase struct {
}

func (a *AuthUsecase) Register(user *entity.User) (*entity.User, error) {
	return user, nil
}
