package handlers

import (
	"encoding/json"
	"fmt"
	"go-postgres-gorm-gin-api/config"
	"go-postgres-gorm-gin-api/models"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"gorm.io/gorm"

	tele "gopkg.in/telebot.v4"
)

type ChatHandler struct {
	DB    *gorm.DB
	TgBot *tele.Bot
}

func NewChatHandler(db *gorm.DB, b *tele.Bot) *ChatHandler {
	return &ChatHandler{DB: db, TgBot: b}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Connections struct {
	MU      sync.Mutex
	Clients map[*websocket.Conn]bool
}

var ActiveConnections = Connections{
	MU:      sync.Mutex{},
	Clients: make(map[*websocket.Conn]bool),
}

func (cn *Connections) AddConnections(c *websocket.Conn) {
	cn.MU.Lock()
	defer cn.MU.Unlock()
	cn.Clients[c] = true
}

func (cn *Connections) RemoveConnections(c *websocket.Conn) {
	cn.MU.Lock()
	defer cn.MU.Unlock()
	c.Close()
	delete(cn.Clients, c)

	fmt.Printf("FROM REMOVE HANDLER: Clients: %v\n", cn.Clients)
}

func (h *ChatHandler) BroadcastMessage(msg *models.Message) {
	ActiveConnections.MU.Lock()
	defer ActiveConnections.MU.Unlock()

	if msg.Platform == models.PlatformWeb {
		if msg.Type == models.MessageChat {

			cfg := config.LoadConfig()
			msgText := fmt.Sprintf("**%s**: %s", msg.Username, msg.Text)

			h.TgBot.Send(&tele.Chat{ID: cfg.TelegramChatId}, msgText, tele.ModeMarkdownV2)
		}
	}

	for client := range ActiveConnections.Clients {
		msgBytes, err := msg.ToJSON()
		if err != nil {
			return
		}

		err = client.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			return
		}

		if err != nil {
			client.Close()
			delete(ActiveConnections.Clients, client)
		}
	}
}

func (h *ChatHandler) ConnectToChatWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	done := make(chan struct{})

	ActiveConnections.AddConnections(conn)

	go func() {
		defer func() {
			ActiveConnections.RemoveConnections(conn)
			done <- struct{}{}
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			var msg models.Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("Error unmarshalling message:", err)
				return
			}

			h.BroadcastMessage(&msg)

		}
	}()

	<-done

}

func (ac *Connections) GetClients() map[*websocket.Conn]bool {
	ac.MU.Lock()
	defer ac.MU.Unlock()
	return ac.Clients
}
