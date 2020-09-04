package mocks

import (
	"peterparada.com/online-bookmarks/domain/entity"
)

type UserRepository struct {
}

func (r *UserRepository) Store(user *entity.User) (error, entity.User) {
	user.ID = "123"
	user.Password = ""

	return nil, *user
}
