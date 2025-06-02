package handlers

import (
	"go-postgres-gorm-gin-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TagHandler struct {
	DB *gorm.DB
}

func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{DB: db}
}

type GetTagsResponse struct {
	Tags []models.Tag `json:"tags"`
}

func (h *TagHandler) GetTags(c *gin.Context) {
	var tags []models.Tag
	if err := h.DB.Find(&tags).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(200, GetTagsResponse{
		Tags: tags,
	})
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.DB.Create(&tag).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(200, tag)
}

type GetTagResponse struct {
	Tag models.Tag `json:"tag"`
}

func (h *TagHandler) GetTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag
	if err := h.DB.First(&tag, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Tag not found"})
		return
	}

	c.JSON(200, GetTagResponse{
		Tag: tag,
	})
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag
	if err := h.DB.First(&tag, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Tag not found"})
		return
	}

	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.DB.Save(&tag).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update tag"})
		return
	}

	c.JSON(200, tag)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag
	if err := h.DB.First(&tag, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Tag not found"})
		return
	}

	if err := h.DB.Delete(&tag).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete tag"})
		return
	}

	c.JSON(200, gin.H{"message": "Tag deleted"})
}
