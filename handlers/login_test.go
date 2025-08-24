package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	jwtSecret = []byte("testsecret")

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "Invalid request body",
			body:       `invalid-json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Wrong credentials",
			body:       `{"username":"user","password":"wrong"}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Correct credentials",
			body:       `{"username":"admin","password":"password"}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			LoginHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK && rr.Body.Len() == 0 {
				t.Errorf("expected response body with token, got empty body")
			}
		})
	}
}
