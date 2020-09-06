package usecase

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
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

func ComparePasswords(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func GenerateAuthToken(claimData map[string]interface{}, jwtSecret string) (string, error) {
	jwtClaims := jwt.MapClaims(claimData)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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

func (a *authUsecase) GenerateAuthToken(userID string, jwtSecret string) (string, error) {
	claimData := map[string]interface{}{
		"id": userID,
	}

	authToken, err := GenerateAuthToken(claimData, jwtSecret)
	if err != nil {
		return "", err
	}

	return authToken, nil
}

func (a *authUsecase) Authenticate(loginData *entity.LoginData, jwtSecret string) (*entity.AuthData, error) {
	user, err := a.userRepo.GetByEmail(loginData.Email)
	if err != nil {
		return nil, err
	}

	passwordMatch := ComparePasswords(user.Password, loginData.Password)
	if passwordMatch == false {
		return nil, errors.New("Invalid password")
	}

	claimData := map[string]interface{}{
		"id": user.ID,
	}

	authToken, err := GenerateAuthToken(claimData, jwtSecret)
	if err != nil {
		return nil, err
	}

	return &entity.AuthData{
		AuthToken: authToken,
	}, nil
}
