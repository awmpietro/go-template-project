package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/nuhorizon/go-project-template/services/template/internal/ports/services"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(jwtService services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			user, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
				return
			}

			// Inject userID into context for handlers
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
