package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateTokens(t *testing.T) {
	token, refreshToken, accessToken := GenerateTokens("user", []byte("WDRwVjdfcU05JXROMXdLNkByRzhqTTJa"), 1*time.Hour, 24*time.Hour)
	fmt.Println(token, refreshToken, accessToken)
}
