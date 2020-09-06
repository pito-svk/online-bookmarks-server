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

func (a *AuthUsecase) RegisterUser(user *entity.User) (*entity.User, error) {
	user.SetID()
	user.SetHashedPassword(user.Password)

	return a.userRepo.Store(user)
}

func (a *AuthUsecase) GenerateAuthToken(userID string, jwtSecret string) (string, error) {
	claimData := map[string]interface{}{
		"id": userID,
	}

	authToken, _ := entity.GenerateAuthToken(claimData, jwtSecret)

	return authToken, nil
}
