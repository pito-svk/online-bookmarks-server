package domain

import "peterparada.com/online-bookmarks/domain/entity"

type AuthUsecase interface {
	RegisterUser(user *entity.User) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GenerateAuthToken(userID string, jwtSecret string) (string, error)
}

type UserRepository interface {
	Store(user *entity.User) (*entity.User, error)
	GetByEmail(email string) (entity.User, error)
}
