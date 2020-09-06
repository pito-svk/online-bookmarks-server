package entity_test

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"peterparada.com/online-bookmarks/domain/entity"
)

func TestGenerateAuthToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		jwtSecret := "SECRET"

		claimData := map[string]interface{}{
			"id": "5f5410bd3cfca9b341bdfe4c",
		}

		authToken, err := entity.GenerateAuthToken(claimData, jwtSecret)

		assert.NoError(t, err)
		assert.NotEmpty(t, authToken)

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		assert.NoError(t, err)

		claims, ok := token.Claims.(jwt.MapClaims)

		assert.True(t, ok)
		assert.NoError(t, claims.Valid())

		assert.Equal(t, claims["id"], "5f5410bd3cfca9b341bdfe4c")
	})
}
