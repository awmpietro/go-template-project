package services

import (
	"errors"

	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
	"github.com/nuhorizon/go-project-template/services/template/internal/ports/services"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) services.JWTService {
	return &jwtService{secretKey: secretKey}
}

func (j *jwtService) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"plan_type": user.PlanType,
		"exp":       domain.TokenExpiry(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(tokenStr string) (*domain.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Recupera o user_id do token
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id missing in token")
	}

	return &domain.User{ID: userID}, nil
}
