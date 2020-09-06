package mocks

import (
	"errors"

	"github.com/google/uuid"
	"peterparada.com/online-bookmarks/auth/usecase"
	"peterparada.com/online-bookmarks/domain/entity"
)

type UserRepository struct {
	users []*entity.User
}

func (repo *UserRepository) Store(user *entity.User) (*entity.User, error) {
	for _, existingUser := range repo.users {
		if existingUser.Email == user.Email {
			return nil, errors.New("Email already exists")
		}
	}

	user.ID = uuid.New().String()

	hashedPassword, err := usecase.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	repo.users = append(repo.users, user)

	return user, nil
}

func (repo *UserRepository) GetByEmail(email string) (entity.User, error) {
	user := entity.User{}

	for _, existingUser := range repo.users {
		if existingUser.Email == email {
			user = entity.User{
				ID:        existingUser.ID,
				Email:     existingUser.Email,
				Password:  existingUser.Password,
				FirstName: existingUser.FirstName,
				LastName:  existingUser.LastName,
			}

			return user, nil
		}
	}

	return user, errors.New("User not found")
}
