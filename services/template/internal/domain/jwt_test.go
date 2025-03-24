package domain

import (
	"os"
	"testing"
	"time"
)

func TestTokenExpiry(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expectedDur time.Duration
	}{
		{"Test with default value", "", 24 * time.Hour},
		{"Test with custom value", "48", 48 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TOKEN_EXPIRE_TIME", tt.envValue)
			} else {
				os.Unsetenv("TOKEN_EXPIRE_TIME")
			}
			defer os.Unsetenv("TOKEN_EXPIRE_TIME")

			now := time.Now().Unix()
			got := TokenExpiry()
			expected := now + int64(tt.expectedDur.Seconds())

			// Aceita uma margem de erro de até 10 segundos por conta do tempo de execução
			if got < expected-10 || got > expected+10 {
				t.Errorf("TokenExpiry() = %v, expected around %v", got, expected)
			}
		})
	}
}
