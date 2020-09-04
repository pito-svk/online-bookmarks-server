package domain

import "peterparada.com/online-bookmarks/domain/entity"

type AuthUsecase interface {
	Register(u *entity.User) (*entity.User, error)
}

type UserRepository interface {
	Store(user *entity.User) (*entity.User, error)
}
