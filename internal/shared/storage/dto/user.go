package dto

type User struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Password string `json:"password" db:"password"`
}
