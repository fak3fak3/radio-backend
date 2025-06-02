package models

import (
	"encoding/json"
	"time"
)

type Platform string

const (
	PlatformTelegram Platform = "telegram"
	PlatformWeb      Platform = "web"
)

type Message struct {
	Text     string    `json:"text"`
	Username string    `json:"username"`
	Date     time.Time `json:"date"`
	Platform Platform  `json:"platform"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
