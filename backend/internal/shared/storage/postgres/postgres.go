package postgres

import (
	"errors"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"service/internal/shared/config"
	"service/internal/shared/storage/dto"
	"service/internal/shared/utils"
)

func InitDB() (*gorm.DB, error) {
	dsn, err := config.GetPostgres()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn.ConnStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&dto.User{}, &dto.Client{})
	if err != nil {
		return nil, err
	}
	err = initializeDefaultClient(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initializeDefaultClient(db *gorm.DB) error {
	var client dto.Client
	result := db.First(&client, "name = ?", "RootClient")
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		jwtSecret, err := utils.GenerateRandomString(32)
		if err != nil {
			return err
		}

		encryptedSecret, err := utils.Encrypt(jwtSecret, config.GetKey())
		if err != nil {
			return err
		}

		defaultClient := dto.Client{
			ID:        0,
			Name:      "RootClient",
			JwtSecret: encryptedSecret,
			Roles:     pq.StringArray{"supreme"},
		}

		if err := db.Create(&defaultClient).Error; err != nil {
			return err
		}

		log.Println("Default client 'RootClient' has been created.")
	} else if result.Error != nil {

		return result.Error
	}

	err := initializeDefaultUser(db, client.ID)
	if err != nil {
		return err
	}

	return nil
}

func initializeDefaultUser(db *gorm.DB, clientID int) error {
	var user dto.User
	result := db.First(&user, "name = ?", "root")
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		passwordHash, err := utils.GenerateHashPassword("your_password_here")
		if err != nil {
			return err
		}

		defaultUser := dto.User{
			ClientID: clientID,
			Name:     "root",
			Password: passwordHash,
			Role:     "supreme",
		}

		if err := db.Create(&defaultUser).Error; err != nil {
			return err
		}

		log.Println("Default user 'root' has been created.")
	} else if result.Error != nil {

		return result.Error
	}

	return nil
}
