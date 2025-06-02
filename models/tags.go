package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model `json:"-"`

	ID    uint    `json:"id" gorm:"primaryKey"`
	Key   string  `json:"key" gorm:"uniqueIndex;not null"`
	Name  string  `json:"name" gorm:"not null"`
	Color *string `json:"color" gorm:"type:char(7);default:null"`
}

type TaggedMedia struct {
	gorm.Model `json:"-"`

	TagID   uint `json:"tag_id" gorm:"primaryKey"`
	MediaID uint `json:"media_id" gorm:"primaryKey"`
}
