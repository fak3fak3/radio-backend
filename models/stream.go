package models

import (
	"time"

	"gorm.io/gorm"
)

type StreamStatus string

const (
	StreamStatusEmpty     StreamStatus = "empty"
	StreamStatusCreated   StreamStatus = "created"
	StreamStatusScheduled StreamStatus = "scheduled"
	StreamStatusRunning   StreamStatus = "running"
	StreamStatusEnded     StreamStatus = "ended"
)

type StreamData struct {
	gorm.Model `json:"-"`

	ID        uint      `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`

	Status    StreamStatus `json:"status" gorm:"type:text"`
	Room      string       `json:"room"`
	TimeStart *time.Time   `json:"time_start"`
	Kbps      int          `json:"kbps"`
	Latency   int64        `json:"latency"`
}
