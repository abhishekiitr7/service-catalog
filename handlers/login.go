package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler creates a JWT token for valid credential.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// For demo, I have used hard-coded credentials.
	if creds.Username != "admin" || creds.Password != "password" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": creds.Username,
		"exp":  time.Now().Add(30 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "could not create token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
