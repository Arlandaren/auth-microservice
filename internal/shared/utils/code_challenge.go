package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

const (
	SHA256 = "SHA256"
	//HMAC   = "HMAC"
)

func generateCodeChallenge(codeVerifier string, methodHash string) (string, error) {

	var codeChallenge string
	if methodHash == SHA256 {
		hash := sha256.New()
		_, err := hash.Write([]byte(codeVerifier))
		if err != nil {
			return "", err
		}
		codeChallenge = base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
	}

	//if methodHash == HMAC {
	//	hash := hmac.New()
	//	_, err := hash.Write([]byte(codeVerifier))
	//	if err != nil {
	//		return "", err
	//	}
	//	codeChallenge = base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
	//}
	return codeChallenge, nil
}

func VerifyCodeVerifier(codeVerifier, methodHash, codeChallengeFromDB string) error {
	codeChallenge, err := generateCodeChallenge(codeVerifier, methodHash)
	if err != nil {
		return err
	}

	if codeChallenge != codeChallengeFromDB {
		return errors.New("не удалось получить CodeChallenge, неверный CodeVerifier")
	}
	return nil
}
