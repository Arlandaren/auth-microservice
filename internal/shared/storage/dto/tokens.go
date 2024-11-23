package dto

import "github.com/golang-jwt/jwt"

type TokenClaims struct {
	Nonce  string
	Scopes []map[string]string
	jwt.StandardClaims
}
