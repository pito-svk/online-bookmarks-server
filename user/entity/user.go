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

func (user *User) ClearPassword() {
	user.Password = ""
}
