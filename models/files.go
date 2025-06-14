package models

import "gorm.io/gorm"

type File struct {
	gorm.Model `json:"-"`

	ID          uint   `json:"id"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Destination string `json:"destination"`

	MediaID *uint `json:"media_id"`
}
