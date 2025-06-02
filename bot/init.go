package bot

import (
	"go-postgres-gorm-gin-api/config"
	"go-postgres-gorm-gin-api/db"
	"go-postgres-gorm-gin-api/handlers"
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

var Bot *tele.Bot

func Init(cfg *config.Config) *tele.Bot {
	db := db.GetPostgresDB()

	var err error

	Bot, err = tele.NewBot(tele.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	botHandler := handlers.NewBotHandler(db)

	Bot.Handle(tele.OnText, func(c tele.Context) error {
		if c.Chat().Type == tele.ChatPrivate {
			return c.Send("This bot is not available in private chats.")
		}

		if c.Chat().Type == tele.ChatGroup || c.Chat().Type == tele.ChatSuperGroup {
			return botHandler.HandleGroupMessage(c)
		}

		return nil
	})

	return Bot
}
