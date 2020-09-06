package domain

import "peterparada.com/online-bookmarks/domain/entity"

type AuthUsecase interface {
	Register(user *entity.User) (*entity.User, error)
	GenerateAuthToken(userID string, jwtSecret string) (string, error)
	Authenticate(loginData *entity.LoginData, jwtSecret string) (*entity.AuthData, error)
}

type UserRepository interface {
	Store(user *entity.User) (*entity.User, error)
	GetByEmail(email string) (entity.User, error)
}
