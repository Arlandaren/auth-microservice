package dto

import "github.com/golang-jwt/jwt"

type TokenIDClaims struct {
	Name   string
	Role   string
	UserID int64
	Nonce  string
	Scopes string
	jwt.StandardClaims
}
