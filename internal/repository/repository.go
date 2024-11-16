package repository

import (
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
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
