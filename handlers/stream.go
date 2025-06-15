package handlers

import (
	"encoding/json"
	"fmt"
	"go-postgres-gorm-gin-api/models"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type StreamHandler struct {
	DB          *gorm.DB
	Connections *Connections
}

func NewStreamHandler(db *gorm.DB) *StreamHandler {
	return &StreamHandler{
		DB: db,
		Connections: &Connections{
			MU:      sync.Mutex{},
			Clients: make(map[*websocket.Conn]bool),
		},
	}
}

func (h *StreamHandler) GetStreamProxy(c *gin.Context) {
	room := c.Query("room")

	url := fmt.Sprintf("http://localhost:7002/live/%s.m3u8", room)

	resp, err := http.Get(url)
	if err != nil {
		c.Status(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		c.Header(k, v[0])
	}
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func (h *StreamHandler) GetStreams(c *gin.Context) {
	var streamData []models.StreamData

	err := h.DB.Find(&streamData).Error
	if err != nil {
		c.JSON(500, &gin.H{"error": err.Error})
	} else {
		c.JSON(200, &gin.H{"streams": streamData})
	}

}

type GetStreamKeyByRoomResponse struct {
	Data string `json:"data"`
}

func (h *StreamHandler) GetStreamKeyByRoom(c *gin.Context) {
	room := c.Query("room")
	url := fmt.Sprintf("http://localhost:8090/control/get?room=%s", room)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(400, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var res GetStreamKeyByRoomResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *StreamHandler) AuthentificateStreamHook(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	log.Println("got hook", string(body))
	log.Println("auth hook called")
	c.JSON(200, gin.H{"code": 0, "msg": "ok"})
}
