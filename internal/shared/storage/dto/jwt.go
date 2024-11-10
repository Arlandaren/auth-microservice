package dto

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
