package utils

import (
	"errors"
	"service/internal/shared/storage/dto"
	"strings"
)

const (
	TokenID     = "TokenID"
	AccessToken = "AccessToken"
)

func CheckScopes(scopes string) []string {
	if scopes == "" {
		return []string{}
	}
	return strings.Split(scopes, " ")
}

func ValidateScopes(scopes []string, typeToken string) error {

	AllowedScopes := []string{"name", "role"}

	if typeToken == TokenID {
		for i := 0; i < len(scopes)+1; i++ {
			if scopes[i] == "role" {
				scopes = append(scopes[:i], scopes[i+1:]...)
			}
		}
	}

	for _, vScopes := range scopes {
		checkScopes := false
		for _, vAllowed := range AllowedScopes {
			if vScopes == vAllowed {
				checkScopes = true
				break // Не нужно проверять дальше
			}
		}
		if !checkScopes {
			return errors.New("данный scope не разрешён для клиента")
		}
	}
	return nil
}

func AddScopes(user *dto.User, scopesFromRequest string, typeToken string) ([]map[string]string, error) {
	scopes := CheckScopes(scopesFromRequest)
	if err := ValidateScopes(scopes, typeToken); err != nil {
		return nil, err
	}

	// Map для хранения значений по ключам
	scopeFromClient := map[string]string{
		"name": user.Name,
		"role": user.Role,
	}

	// Слайс мап
	var mapScopesSlice []map[string]string

	// Перебираем запрошенные области и создаём мапу
	for _, scope := range scopes {
		if value, exists := scopeFromClient[scope]; exists {
			mapScopesSlice = append(mapScopesSlice, map[string]string{
				scope: value,
			})
		}
	}

	return mapScopesSlice, nil
}
