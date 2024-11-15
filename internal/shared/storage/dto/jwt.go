package dto

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
