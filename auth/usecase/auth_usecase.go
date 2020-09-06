package usecase

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

func (authU *authUsecase) Register(u *entity.User) (*entity.User, error) {
	u.SetID()

	err := u.SetHashedPassword(u.Password)
	if err != nil {
		return nil, err
	}

	user, err := authU.userRepo.Store(u)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (a *authUsecase) GenerateAuthToken(userID string, jwtSecret string) (string, error) {
	claimData := map[string]interface{}{
		"id": userID,
	}

	authToken, err := entity.GenerateAuthToken(claimData, jwtSecret)
	if err != nil {
		return "", err
	}

	return authToken, nil
}
