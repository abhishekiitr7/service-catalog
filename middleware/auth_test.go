package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestAuthMiddleware(t *testing.T) {
	secret := []byte("test_secret_key_32_characters_long")
	SetJWTSecretForTesting(secret)
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": "abc"})
	tokenString, _ := token.SignedString(secret)

	tests := []struct {
		name       string
		path       string
		authHeader string
		wantCode   int
	}{
		{
			name:       "Allow public login",
			path:       "/login",
			authHeader: "",
			wantCode:   http.StatusOK,
		},
		{
			name:       "No Authorization header",
			path:       "/protected",
			authHeader: "",
			wantCode:   http.StatusUnauthorized,
		},
		{
			name:       "Invalid token format",
			path:       "/protected",
			authHeader: "InvalidToken",
			wantCode:   http.StatusUnauthorized,
		},
		{
			name:       "Valid token",
			path:       "/protected",
			authHeader: "Bearer " + tokenString,
			wantCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := AuthMiddleware(next)

			req := httptest.NewRequest("GET", tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantCode)
			}
		})
	}
}
