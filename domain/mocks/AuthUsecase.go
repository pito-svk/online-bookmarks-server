package mocks

import (
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type authUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthUsecase(userRepo domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{
		userRepo,
	}
}

func (a *authUsecase) RegisterUser(user *entity.User) (*entity.User, error) {
	user.SetID()
	user.SetHashedPassword()

	return a.userRepo.Store(user)
}

func (a *authUsecase) GenerateAuthToken(userID string, jwtSecret string) (string, error) {
	claimData := map[string]interface{}{
		"id": userID,
	}

	authToken, _ := entity.GenerateAuthToken(claimData, jwtSecret)

	return authToken, nil
}
