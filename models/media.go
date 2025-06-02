package models

import "gorm.io/gorm"

type MediaType string

const (
	MediaTypeAudioSelfHosted MediaType = "audio_self_hosted"
	MediaTypeAudioSoundCloud MediaType = "audio_soundcloud"
	MediaTypeVideoSelfHosted MediaType = "video_self_hosted"
	MediaTypeVideoYouTube    MediaType = "video_youtube"
)

type Media struct {
	gorm.Model

	ID          int       `json:"id" gorm:"primaryKey"`
	Type        MediaType `json:"type" gorm:"type:text"`
	Duration    int       `json:"duration"`
	Cover       *File     `json:"cover" gorm:"foreignKey:MediaID;references:ID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Url         string    `json:"url"`

	Tags []*Tag `json:"tags" gorm:"many2many:tagged_media;"`
}
