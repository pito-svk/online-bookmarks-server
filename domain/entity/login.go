package entity

import "github.com/dgrijalva/jwt-go"

type LoginData struct {
	Email    string
	Password string
}

type AuthData struct {
	Token string
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
