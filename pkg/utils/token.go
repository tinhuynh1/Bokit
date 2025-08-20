package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	Role    string
	TokenID string
	jwt.RegisteredClaims
}

func GenerateTokens(role string, jwtSecret []byte, accessTTL, refreshTTL time.Duration) (string, string, string) {
	tokenID := uuid.NewString()

	accessClaims := CustomClaims{
		Role:    role,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
		},
	}

	refreshClaims := CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTTL)),
		},
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, _ := access.SignedString(jwtSecret)
	refreshToken, _ := refresh.SignedString(jwtSecret)

	return accessToken, refreshToken, tokenID
}

func VerifyAccessToken(tokenStr string, jwtSecret []byte) (*CustomClaims, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token is empty")
	}

	// Trim whitespace
	tokenStr = strings.TrimSpace(tokenStr)

	// Kiểm tra format cơ bản
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("token has %d parts, expected 3. Token: %s", len(parts), tokenStr)
	}

	// Kiểm tra header có thể decode được không
	_, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid header base64: %v. Header part: %s", err, parts[0])
	}

	// Log để debug
	fmt.Printf("Token parts: %d, Header: %s, Payload: %s, Signature: %s\n",
		len(parts), parts[0], parts[1], parts[2])

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(*CustomClaims), nil
}

func VerifyRefreshToken(tokenStr string, jwtSecret []byte) (*CustomClaims, error) {
	return VerifyAccessToken(tokenStr, jwtSecret)
}
