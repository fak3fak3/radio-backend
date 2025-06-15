package handlers

import (
	"fmt"
	"go-postgres-gorm-gin-api/models"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (h *StreamHandler) GetStreamCredentials(c *gin.Context) {
	var credentials []models.StreamCredentials

	err := h.DB.Find(&credentials).Error
	if err != nil {
		c.JSON(500, &gin.H{"error": err.Error})
	} else {
		c.JSON(200, &gin.H{"credentials": credentials})
	}
}

func (h *StreamHandler) DeleteStreamCredentials(c *gin.Context) {
	id := c.Param("id")
	var creds models.StreamCredentials
	if err := h.DB.First(&creds, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Credentails not found"})
		return
	}

	if err := h.DB.Delete(&creds).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete credentials"})
		return
	}

	var credentials []models.StreamCredentials
	if err := h.DB.Find(&credentials).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"credentials": credentials})
}

type CreateStreamCredentialsRequest struct {
	Room     string `json:"room"`
	Password string `json:"password"`
}

func (h *StreamHandler) CreateStreamCredentials(c *gin.Context) {
	var req CreateStreamCredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	creds := models.StreamCredentials{
		Room:     req.Room,
		Password: req.Password,
	}

	err := h.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "room"}}, // уникальный ключ
		DoUpdates: clause.AssignmentColumns([]string{"password", "updated_at"}),
	}).Create(&creds).Error

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var credentials []models.StreamCredentials
	if err := h.DB.Find(&credentials).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"credentials": credentials})
}
