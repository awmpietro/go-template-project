package services

import (
	"github.com/nuhorizon/go-project-template/services/template/internal/domain"
)

type JWTService interface {
	GenerateToken(user *domain.User) (string, error)
	ValidateToken(token string) (*domain.User, error)
}
