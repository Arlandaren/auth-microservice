package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"service/internal/shared/config"
	"service/internal/shared/storage/dto"
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

	err = db.AutoMigrate(
		&dto.User{},
		&dto.AuthCodeOIDC{},
		&dto.ClientData{},
		&dto.AccessToken{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
