package utils

import (
	"github.com/dgrijalva/jwt-go"
	"service/internal/shared/config"
	"service/internal/shared/storage/dto"
	"time"
)

func GenerateToken(userId int) (string, error) {
	jwtKey := []byte(config.GetJwt())
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &dto.Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Auth",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
