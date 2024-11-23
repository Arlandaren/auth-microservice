package dto

import "time"

type AuthCodeOIDC struct {
	AuthCode            string    `json:"auth_code" db:"auth_code"`
	UserID              int       `json:"user_id" db:"user_id"`
	ClientID            string    `json:"client_id" db:"client_id"`
	RedirectURI         string    `json:"redirect_uri" db:"redirect_uri"`
	Scopes              string    `json:"scopes" db:"scopes"`
	State               string    `json:"state" db:"state"`
	CodeChallenge       string    `json:"code_challenge" db:"code_challenge"`
	CodeChallengeMethod string    `json:"code_challenge_method" db:"code_challenge_method"`
	ExpiresAt           time.Time `json:"expires_at" db:"expires_at"`
}
