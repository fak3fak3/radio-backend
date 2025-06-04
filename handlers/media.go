package handlers

import (
	"go-postgres-gorm-gin-api/models"
	"go-postgres-gorm-gin-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/davecgh/go-spew/spew"
)

type MediaHandler struct {
	DB *gorm.DB
}

func NewMediaHandler(db *gorm.DB) *MediaHandler {
	return &MediaHandler{DB: db}
}

type GetMediaResponse struct {
	Media []models.Media `json:"media"`
}

func (h *MediaHandler) GetResults(c *gin.Context) {
	var media []models.Media
	if err := h.DB.Preload("Tags").Preload("Cover").Find(&media).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch media"})
		return
	}

	for i := range media {
		media[i].Cover.Path = utils.GenerateFullFilePath(media[i].Cover.Path)
	}

	c.JSON(200, GetMediaResponse{
		Media: media,
	})
}

func (h *MediaHandler) CreateMedia(c *gin.Context) {
	var media models.Media

	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	spew.Dump(media)

	if err := h.DB.Create(&media).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(200, "ok")
}

func (h *MediaHandler) UpdateMediaTags(m *models.Media, ts []*models.Tag) error {
	if err := h.DB.Model(m).Association("Tags").Replace(m.Tags); err != nil {
		return err
	}
	return nil
}
