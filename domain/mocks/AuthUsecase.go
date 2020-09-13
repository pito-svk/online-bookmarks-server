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
	err := user.SetHashedPassword()
	if err != nil {
		panic(err)
	}

	return a.userRepo.Store(user)
}

func (a *authUsecase) GenerateAuthToken(userID string, jwtSecret string) (string, error) {
	claimData := map[string]interface{}{
		"id": userID,
	}

	authToken, err := entity.GenerateAuthToken(claimData, jwtSecret)
	if err != nil {
		panic(err)
	}

	return authToken, nil
}
