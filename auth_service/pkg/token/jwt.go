package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateNewToken() *jwt.Token {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(43000, 0)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}
