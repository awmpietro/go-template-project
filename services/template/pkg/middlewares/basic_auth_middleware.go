package middlewares

import (
	"net/http"
	"os"
)

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		expectedUser := os.Getenv("SWAGGER_USER_AUTH")
		expectedPass := os.Getenv("SWAGGER_PASSWORD_AUTH")

		if !ok || user != expectedUser || pass != expectedPass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
