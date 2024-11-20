package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"service/internal/shared/config"
	"service/internal/shared/storage/dto"
	"time"
)

func GenerateToken(userId int, role string) (string, error) {
	jwtKey := []byte(config.GetJwt())
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &dto.Claims{
		UserID: userId,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth_service",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*dto.Claims, error) {
	secretKey := []byte(config.GetJwt())

	claims := &dto.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
