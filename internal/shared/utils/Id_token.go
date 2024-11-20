package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/thanhpk/randstr"
	"log"
	"service/internal/shared/config"
	"service/internal/shared/storage/dto"
	"time"
)

func GenerateKeysRSA() (Private []byte, Public []byte) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Не удалось сгенерировать приватный RSA ключ: %v", err)
	}

	publicBytes, err := x509.MarshalPKCS8PrivateKey(privateKey.PublicKey)
	if err != nil {
		log.Fatalf("Не удалось маршализовать публичный ключ: %v", err)
	}
	pemKeyPBL := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicBytes})

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Не удалось маршализовать приватный ключ: %v", err)
	}
	pemKeyPR := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})

	return pemKeyPR, pemKeyPBL
}

func scopesToString(scopes []string) string {
	return fmt.Sprintf("%s", scopes)
}

func GenerateIDToken(userID, clientID, issuer string) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(config.GetPrivateKey())
	if err != nil {
		return "", err
	}

	claims := dto.TokenIDClaims{
		Name:   "",
		Role:   "",
		Nonce:  randstr.String(10),
		Scopes: scopesToString(scopes),
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  clientID,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// jwt token шифруется RSA алгоритмом
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
