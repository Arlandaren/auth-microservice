package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
)

func GenerateAuthCode(userID int) (string, error) {
	// Создаем массив случайных байтов
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации случайных данных: %v", err)
	}

	// Кодируем в base64 (или можно использовать hex)
	authCode := base64.URLEncoding.EncodeToString(b)

	return fmt.Sprintf("%s.%s", authCode, strconv.Itoa(userID)), nil
}
