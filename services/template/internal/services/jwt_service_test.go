package services_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
	internalservices "github.com/nuhorizon/go-project-template/services/template/internal/services"
	"github.com/stretchr/testify/assert"
)

const testSecretKey = "supersecretkey"

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	jwtService := internalservices.NewJWTService(testSecretKey)

	t.Run("successfully generates and validates token", func(t *testing.T) {
		user := &domain.User{ID: "user-123", PlanType: "premium"}

		tokenStr, err := jwtService.GenerateToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenStr)

		validatedUser, err := jwtService.ValidateToken(tokenStr)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, validatedUser.ID)
	})

	t.Run("fails validation with invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		user, err := jwtService.ValidateToken(invalidToken)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("fails validation if user_id missing in token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"plan_type": "premium",
			"exp":       time.Now().Add(1 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, _ := token.SignedString([]byte(testSecretKey))

		user, err := jwtService.ValidateToken(tokenStr)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
