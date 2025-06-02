package handlers

import (
	"go-postgres-gorm-gin-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	if err := h.DB.Preload("Tags").Find(&media).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch media"})
		return
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

	if media.Type == models.MediaTypeAudioSoundCloud {
		media, err := h.CreateSoundcloudMedia(&media)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create media"})
			return
		}

		c.JSON(200, media)
	}
}

func (h *MediaHandler) CreateSoundcloudMedia(m *models.Media) (*models.Media, error) {
	if err := h.DB.Create(&m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (h *MediaHandler) UpdateMediaTags(m *models.Media, ts []*models.Tag) error {
	if err := h.DB.Model(m).Association("Tags").Replace(m.Tags); err != nil {
		return err
	}
	return nil
}
