package repository

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"service/internal/shared/storage/dto"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(wrapper *gorm.DB) *Repository {
	log.Println("NewRepository")
	return &Repository{
		db: wrapper,
	}
}

func (r *Repository) NewUser(user *dto.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) NewAuthCode(authCode *dto.AuthCodeOIDC) error {
	result := r.db.Create(authCode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) NewClientData(clientData *dto.ClientData) error {
	result := r.db.Create(clientData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) NewAccessToken(accessToken *dto.AccessToken) error {
	result := r.db.Create(accessToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) GetUserById(userId int) (*dto.User, error) {
	var user dto.User
	result := r.db.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *Repository) GetUserByName(userName string) (*dto.User, error) {
	var user dto.User
	result := r.db.Where("name = ?", userName).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *Repository) CheckGetClientID(clientID string) error {
	var clientData dto.ClientData
	result := r.db.Where("client_id = ?", clientID).First(&clientData)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("ClientID %s не найден.\n", clientID)
		} else {
			return fmt.Errorf("Ошибка при запросе: %s\n", result.Error)
		}
	}

	return nil
}

// Unused
//func (r *Repository) CheckClientSecret(clientSecret string) (bool, error) {
//	return true, nil
//}

func (r *Repository) GetAuthCodeFromClientID(clientID string) (*dto.AuthCodeOIDC, error) {
	var authCodeOIDC dto.AuthCodeOIDC
	result := r.db.Where("client_id = ?", clientID).First(&authCodeOIDC)
	if result.Error != nil {
		return nil, result.Error
	}
	return &authCodeOIDC, nil
}

func (r *Repository) GetClientIDandClientSecret(clientID string) (*dto.ClientData, error) {
	var clientData dto.ClientData
	result := r.db.Where("client_id = ?", clientID).First(&clientData)
	if result.Error != nil {
		return nil, result.Error
	}
	return &clientData, nil
}

// Unused
//func (r *Repository) AddClientIDandClientSecret(clientID, ClientSecret string) error {
//	return nil
//}

func (r *Repository) DeleteAuthCodeFromClientID(clientID string) error {
	var authCodeOIDC dto.AuthCodeOIDC
	result := r.db.Where("client_id = ?", clientID).Delete(&authCodeOIDC)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
