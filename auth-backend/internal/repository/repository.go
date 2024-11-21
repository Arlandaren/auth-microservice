package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"service/internal/shared/config"
	"service/internal/shared/utils"
	"time"

	log "github.com/sirupsen/logrus"
	"service/internal/shared/storage/dto"
)

type Repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRepository(wrapper *gorm.DB, redis *redis.Client) *Repository {
	log.Println("NewRepository")
	return &Repository{
		db:  wrapper,
		rdb: redis,
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

func (r *Repository) GetUserByNameForClient(userName string, clientId int) (*dto.User, error) {
	var user dto.User
	result := r.db.Where("name = ? AND client_id = ?", userName, clientId).First(&user)
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

func (r *Repository) GetClientByName(name string) (*dto.Client, error) {
	var client dto.Client
	result := r.db.Where("name = ?", name).First(&client)
	if result.Error != nil {
		return nil, result.Error
	}
	return &client, nil
}

func (r *Repository) NewClient(client *dto.Client) error {
	result := r.db.Create(client)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) GetClientByID(clientID int, ctx context.Context) (*dto.Client, error) {
	var client dto.Client

	redisKey := fmt.Sprintf("client:%d", clientID)
	cachedClientData, err := r.rdb.Get(ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		result := r.db.Where("id = ?", clientID).First(&client)
		if result.Error != nil {
			return nil, result.Error
		}

		clientData, err := json.Marshal(&client)
		if err != nil {
			return nil, err
		}

		err = r.rdb.Set(ctx, redisKey, clientData, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Ошибка при сохранении клиента в Redis: %v", err)
		}
	} else if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal([]byte(cachedClientData), &client)
		if err != nil {
			return nil, err
		}
	}

	return &client, nil
}

func (r *Repository) GetClientJwtSecret(clientID int, ctx context.Context) (string, error) {
	var jwtSecret string

	redisKey := fmt.Sprintf("client_jwt_secret:%d", clientID)

	cachedJwtSecret, err := r.rdb.Get(ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {

		var client dto.Client
		result := r.db.Select("jwt_secret").Where("id = ?", clientID).First(&client)
		if result.Error != nil {
			return "", result.Error
		}

		jwtSecret = client.JwtSecret

		err = r.rdb.Set(ctx, redisKey, jwtSecret, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Ошибка при сохранении jwt_secret в Redis: %v", err)
		}
	} else if err != nil {
		return "", err
	} else {
		jwtSecret = cachedJwtSecret
	}

	decryptedJwtSecret, err := utils.Decrypt(jwtSecret, config.GetKey())
	if err != nil {
		return "", err
	}

	jwtSecret = decryptedJwtSecret

	return jwtSecret, nil
}
