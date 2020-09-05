package filedb

import (
	"errors"
	"log"

	"github.com/nanobox-io/golang-scribble"
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type fileDBUserRepo struct {
	DB *scribble.Driver
}

func NewFileDBUserRepository(db *scribble.Driver) domain.UserRepository {
	return &fileDBUserRepo{
		DB: db,
	}
}

func (userRepo *fileDBUserRepo) Store(user *entity.User) (*entity.User, error) {
	err := userRepo.DB.Read("userdata", user.Email, user)
	if err == nil {
		return nil, errors.New("User already exists")
	}

	err = userRepo.DB.Write("userdata", user.Email, user)
	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}