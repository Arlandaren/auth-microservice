package dto

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserID   int    `json:"user_id"`
	Role     string `json:"role"`
	ClientID int    `json:"client_id"`
	jwt.StandardClaims
}