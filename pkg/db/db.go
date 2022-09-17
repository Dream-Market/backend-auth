package db

import (
	"fmt"
	"log"
	"time"

	"github.com/Dream-Market/backend-auth/pkg/config"
	"github.com/Dream-Market/backend-auth/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB             *gorm.DB
	UserHandler    *UserHandler
	SessionHandler *SessionHandler
}

func Init(c config.Config) Handler {
	var err error
	var db *gorm.DB
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	for i := 0; i < int(c.DBConnectionRetries); i++ {
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
		if err != nil {
			time.Sleep(time.Duration(int64(time.Second) * c.DBConnectionInterval))
		} else {
			break
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.User{}, &models.Session{})

	h := Handler{
		DB: db,
	}
	h.UserHandler = InitUserHandler(&h)
	h.SessionHandler = InitSessionHandler(&h, c.ExpirationHours)

	return h
}
