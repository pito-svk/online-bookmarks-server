package usecase

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (a *authUsecase) Register(u *entity.User) (*entity.User, error) {
	u.ID = GenerateID()
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = hashedPassword

	user, err := a.userRepo.Store(u)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}
