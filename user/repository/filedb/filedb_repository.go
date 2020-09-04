package filedb

import (
	"github.com/boltdb/bolt"
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type fileDBUserRepo struct {
	DB *bolt.DB
}

func NewFileDBUserRepository(db *bolt.DB) domain.UserRepository {
	return &fileDBUserRepo{
		DB: db,
	}
}

func (userRepo *fileDBUserRepo) Store(user *entity.User) (*entity.User, error) {
	// TODO: Implement
	return nil, nil
}
