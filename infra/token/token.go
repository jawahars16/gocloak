package token

import "github.com/golang-jwt/jwt/v5"

type generator struct{}

func NewGenerator() generator {
	return generator{}
}

func (g *generator) NewWithClaims(method jwt.SigningMethod, claims jwt.Claims, opts ...jwt.TokenOption) *jwt.Token {
	return jwt.NewWithClaims(method, claims, opts...)
}
