package mocks

import (
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type AuthUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthUsecase(userRepo domain.UserRepository) domain.AuthUsecase {
	return &AuthUsecase{
		userRepo,
	}
}

func (a *AuthUsecase) Register(user *entity.User) (*entity.User, error) {
	user.ID = entity.GenerateID()
	hashedPassword, _ := entity.HashPassword(user.Password)

	user.Password = hashedPassword

	return a.userRepo.Store(user)
}

func (a *AuthUsecase) GenerateAuthToken(userID string, jwtSecret string) (string, error) {
	claimData := map[string]interface{}{
		"id": userID,
	}

	authToken, _ := entity.GenerateAuthToken(claimData, jwtSecret)

	return authToken, nil
}
