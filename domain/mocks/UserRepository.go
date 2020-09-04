package mocks

import (
	"github.com/google/uuid"
	"peterparada.com/online-bookmarks/domain/entity"
)

type UserRepository struct {
}

func (r *UserRepository) Store(user *entity.User) (error, entity.User) {
	user.ID = uuid.New().String()
	user.Password = ""

	return nil, *user
}
