package dto

import (
	"github.com/lib/pq"
)

type Client struct {
	ID        int            `gorm:"primaryKey;column:id"`
	Name      string         `gorm:"column:name"`
	JwtSecret string         `gorm:"column:jwt_secret"`
	Roles     pq.StringArray `gorm:"type:text[];column:roles"`
}
