package handlers

import (
	"go-postgres-gorm-gin-api/models"

	"github.com/gorilla/websocket"
	tele "gopkg.in/telebot.v4"
	"gorm.io/gorm"
)

type BotHandler struct {
	DB *gorm.DB
}

func NewBotHandler(db *gorm.DB) *BotHandler {
	return &BotHandler{DB: db}
}

func (h *BotHandler) HandleGroupMessage(c tele.Context) error {
	cs := ActiveConnections.GetClients()

	msg := models.Message{
		Text:     c.Text(),
		Username: c.Sender().Username,
		Date:     c.Message().Time(),
		Platform: models.PlatformTelegram,
	}

	msgBytes, err := msg.ToJSON()
	if err != nil {
		return err
	}

	for client := range cs {
		err := client.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			return err
		}
	}

	return c.Reply(c.Text)
}
