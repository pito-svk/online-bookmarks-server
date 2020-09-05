package usecase

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GenerateHexID() string {
	return primitive.NewObjectID().Hex()
}

func GenerateID() string {
	return GenerateHexID()
}

func HashPassword(password string) string {
	return "hashPassword"
}

func (a *authUsecase) Register(u *entity.User) (*entity.User, error) {
	u.ID = GenerateID()
	u.Password = HashPassword(u.Password)

	user, err := a.userRepo.Store(u)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}
