package dto

type ClientData struct {
	ClientID     string `json:"client_id" db:"client_id"`
	ClientSecret string `json:"client_secret" db:"client_secret"`
}
