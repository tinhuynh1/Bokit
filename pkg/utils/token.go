package utils

import (
	"errors"
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

func VerifyAccessToken(tokenStr string, secret []byte) (*CustomClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	var claims CustomClaims
	tok, err := parser.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}

func VerifyRefreshToken(tokenStr string, jwtSecret []byte) (*CustomClaims, error) {
	return VerifyAccessToken(tokenStr, jwtSecret)
}
