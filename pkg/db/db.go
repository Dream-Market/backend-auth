package db

import (
	"log"

	"github.com/Dream-Market/backend-auth/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB             *gorm.DB
	UserHandler    *UserHandler
	SessionHandler *SessionHandler
}

func Init(url string, expirationHours int64) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.User{}, &models.Session{})

	h := Handler{
		DB: db,
	}
	h.UserHandler = InitUserHandler(&h)
	h.SessionHandler = InitSessionHandler(&h, expirationHours)

	return h
}
