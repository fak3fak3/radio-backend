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

func (h *ChatHandler) AddConnections(c *websocket.Conn) {
	ActiveConnections.MU.Lock()
	defer ActiveConnections.MU.Unlock()
	ActiveConnections.Clients[c] = true
}

func (h *ChatHandler) RemoveConnections(c *websocket.Conn) {
	ActiveConnections.MU.Lock()
	defer ActiveConnections.MU.Unlock()
	c.Close()
	delete(ActiveConnections.Clients, c)

	fmt.Printf("FROM REMOVE HANDLER: Clients: %v\n", ActiveConnections.Clients)
}

func (h *ChatHandler) BroadcastMessage(message []byte) {
	ActiveConnections.MU.Lock()
	defer ActiveConnections.MU.Unlock()

	var msg models.Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Println("Error unmarshalling message:", err)
		return
	}

	if msg.Platform == models.PlatformWeb {
		cfg := config.LoadConfig()
		msgText := fmt.Sprintf("**%s**: %s", msg.Username, msg.Text)

		h.TgBot.Send(&tele.Chat{ID: cfg.TelegramChatId}, msgText, tele.ModeMarkdownV2)
	}

	for client := range ActiveConnections.Clients {
		err := client.WriteMessage(websocket.TextMessage, message)

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

	h.AddConnections(conn)

	go func() {
		defer func() {
			h.RemoveConnections(conn)
			done <- struct{}{}
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			h.BroadcastMessage(message)

		}
	}()

	<-done

}

func (ac *Connections) GetClients() map[*websocket.Conn]bool {
	ac.MU.Lock()
	defer ac.MU.Unlock()
	return ac.Clients
}
