package dto

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	ClientID int    `json:"client_id" gorm:"not null;index"`
	Client   Client `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
