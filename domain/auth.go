package domain

import "peterparada.com/online-bookmarks/domain/entity"

type AuthUsecase interface {
	Register(u *entity.User) (error, entity.User)
}

type UserRepository interface {
	Store(user *entity.User) (error, entity.User)
}
