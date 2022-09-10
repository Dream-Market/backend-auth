package models

import "time"

type Session struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	UserId    int64     `json:"user_id,omitempty"`
	IssuedAt  time.Time `gorm:"issued_at"`
	ExpiresAt time.Time `gorm:"expires_at"`
	IsBlocked bool      `gorm:"is_blocked"`
}
