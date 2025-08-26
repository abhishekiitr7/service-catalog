package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	validateJWTSecret()
}

func validateJWTSecret() {
	if len(jwtSecret) == 0 {
		panic("JWT_SECRET environment variable is required")
	}
	if len(jwtSecret) < 32 {
		panic("JWT_SECRET must be at least 32 characters long for security")
	}
}

// SetJWTSecretForTesting allows tests to override the JWT secret
func SetJWTSecretForTesting(secret []byte) {
	jwtSecret = secret
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow public login
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
