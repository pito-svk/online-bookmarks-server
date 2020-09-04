package filedb

import (
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type fileDBUserRepo struct {
}

func NewFileDBUserRepository() domain.UserRepository {
	return &fileDBUserRepo{}
}

func (userRepo *fileDBUserRepo) Store(user *entity.User) (*entity.User, error) {
	// TODO: Implement
	return nil, nil
}
