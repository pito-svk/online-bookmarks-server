package entity

import "github.com/dgrijalva/jwt-go"

// TODO: Shouldn't be part of entity

func GenerateAuthToken(claimData map[string]interface{}, jwtSecret string) (string, error) {
	jwtClaims := jwt.MapClaims(claimData)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
