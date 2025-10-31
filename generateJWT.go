package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	DefaultLoginTokenDuration = 24 * time.Hour
)

type CustomClaims struct {
	Email string `json:"email"`

	jwt.RegisteredClaims
}

type JWTGenerator struct {
	secret        []byte
	signingMethod jwt.SigningMethod
}

func NewJWTGenerator(secret []byte) *JWTGenerator {
	return &JWTGenerator{
		secret:        secret,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (g *JWTGenerator) generateToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(g.signingMethod, claims)
	tokenStr, err := token.SignedString(g.secret)
	if err != nil {
		return "", fmt.Errorf("falha ao assinar token JWT com o método %s: %w", g.signingMethod.Alg(), err)
	}
	return tokenStr, nil
}

func (g *JWTGenerator) GenerateLoginToken(email string) (string, error) {
	expirationTime := time.Now().Add(DefaultLoginTokenDuration)
	claims := &CustomClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return g.generateToken(claims)
}
