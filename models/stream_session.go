package models

import (
	"time"

	"gorm.io/gorm"
)

type StreamCredentials struct {
	gorm.Model `json:"-"`

	ID       uint   `json:"id"`
	Room     string `json:"room"`
	Password string `json:"password"`

	Sessions []StreamSession `gorm:"foreignKey:StreamCredentialsID;references:ID"`
}

type StreamSession struct {
	gorm.Model `json:"-"`

	CreatedAt time.Time `json:"created_at"`

	StreamCredentialsID uint               `json:"stream_credentials_id"`
	StreamCredentials   *StreamCredentials `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
