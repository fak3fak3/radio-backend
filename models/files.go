package models

import "gorm.io/gorm"

type File struct {
	gorm.Model `json:"-"`

	ID   uint   `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`

	MediaID *uint `json:"media_id"`
}
