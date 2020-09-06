package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Email     string
	Password  string
	FirstName string
	LastName  string
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

func (user *User) SetID() {
	user.ID = GenerateID()
}

func (user *User) SetHashedPassword(password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return nil
}

func (user *User) ClearPassword() {
	user.Password = ""
}
