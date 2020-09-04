package mocks

import (
	"errors"

	"github.com/google/uuid"
	"peterparada.com/online-bookmarks/domain/entity"
)

type UserRepository struct {
	users []*entity.User
}

func (repo *UserRepository) Store(user *entity.User) (*entity.User, error) {
	for _, existingUser := range repo.users {
		if existingUser.Email == user.Email {
			return nil, errors.New("User already exists")
		}
	}

	user.ID = uuid.New().String()
	user.Password = ""

	repo.users = append(repo.users, user)

	return user, nil
}
