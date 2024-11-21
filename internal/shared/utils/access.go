package utils

import "strings"

func HasAccess(role string, allowedRoles ...string) bool {
	role = strings.TrimSpace(role)
	for _, allowedRole := range allowedRoles {
		if role == strings.TrimSpace(allowedRole) {
			return true
		}
	}
	return false
}
