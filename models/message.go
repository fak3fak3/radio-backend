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

type MessageType string

const (
	MessageChat MessageType = "chat"
	MessageInfo MessageType = "info"
)

type Message struct {
	Text     string      `json:"text"`
	Username string      `json:"username"`
	Date     time.Time   `json:"date"`
	Platform Platform    `json:"platform"`
	Type     MessageType `json:"type"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
