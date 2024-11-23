package dto

import "time"

type AccessToken struct {
	Token     string
	ClientID  string
	UserID    int
	TokenType string
	ExpiresAt time.Time
}
