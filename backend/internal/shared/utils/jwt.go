package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"service/internal/shared/storage/dto"
	"time"
)

func GenerateToken(userId, clientId int, role string, secretKey []byte) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &dto.Claims{
		UserID:   userId,
		Role:     role,
		ClientID: clientId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth_service",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string, secretKey []byte) (*dto.Claims, error) {

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

func ParseTokenWithoutVerification(tokenStr string) (*dto.Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &dto.Claims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*dto.Claims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
