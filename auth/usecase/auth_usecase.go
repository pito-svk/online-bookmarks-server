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

type userDataInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func generateHexID() string {
	return primitive.NewObjectID().Hex()
}

func generateID() string {
	return generateHexID()
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (user *User) SetID() {
	user.ID = generateID()
}

func (user *User) SetHashedPassword() error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return nil
}

func (authU *authUsecase) RegisterUser(u userData) (*entity.User, error) {
	u.SetID()

	err := u.SetHashedPassword()
	if err != nil {
		return nil, err
	}

	userDataInputForEntity := formatWithId(u)

	user, err := authU.userRepo.Store(userDataInputForEntity)
	if err != nil {
		return nil, err
	}

	user.ClearPassword()

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
