package utils

import (
	"github.com/golang-jwt/jwt"
	"github.com/thanhpk/randstr"
	"os"
	"path/filepath"
	"service/internal/shared/storage/dto"
	"time"
)

const (
	PrivateKeyAccessToken = "AccessKey"
	PrivateKeyTokenID     = "IDKey"
)

func GenerateTokenOIDC(userID, clientID, issuer string, scopes []map[string]string, typeTokenKey string) (string, error) {
	var getPrivateKey []byte
	var err error
	if typeTokenKey == PrivateKeyAccessToken {
		getPrivateKey, err = os.ReadFile(filepath.Join("utils", "rsaAccessTokenPrivateKey.pem"))
		if err != nil {
			return "", err
		}
	} else if typeTokenKey == PrivateKeyTokenID {
		getPrivateKey, err = os.ReadFile(filepath.Join("utils", "rsaTokenIDPrivateKey.pem"))
		if err != nil {
			return "", err
		}
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(getPrivateKey)
	if err != nil {
		return "", err
	}

	claims := dto.TokenClaims{
		Nonce:  randstr.String(10),
		Scopes: scopes,
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  clientID,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// jwt token подписывается RSA алгоритмом
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
