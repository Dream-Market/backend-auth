package db

import (
	"time"

	"github.com/Dream-Market/backend-auth/pkg/models"
)

type SessionHandler struct {
	*Handler
	ExpirationHours time.Duration
}

func InitSessionHandler(h *Handler, expirationHours int64) *SessionHandler {
	return &SessionHandler{
		Handler:         h,
		ExpirationHours: time.Duration(expirationHours),
	}
}

func (h *SessionHandler) Create(s models.Session) (result models.Session, err error) {
	s.ExpiresAt = time.Now().Add(h.ExpirationHours)
	err = h.DB.Create(&s).Error
	return s, err
}

func (h *SessionHandler) IsActive(id int64) (active bool, err error) {
	var session models.Session
	err = h.DB.Where(&models.Session{Id: id}).First(&session).Error
	if err != nil {
		return
	}

	if session.IsBlocked {
		return
	}
	if session.ExpiresAt.Before(time.Now()) {
		h.Block(id)
		return
	}
	return true, err
}

func (h *SessionHandler) Block(id int64) (err error) {
	return h.DB.Where(&models.Session{Id: id}).Updates(&models.Session{IsBlocked: true}).Error
}
