package handlers

import (
	"go-postgres-gorm-gin-api/config"
	"go-postgres-gorm-gin-api/models"
	"go-postgres-gorm-gin-api/s3"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileHandler struct {
	DB *gorm.DB
	S3 *s3.S3Client
}

func NewFileHandler(db *gorm.DB, s3 *s3.S3Client) *FileHandler {
	return &FileHandler{
		DB: db,
		S3: s3,
	}
}

type FileUploadResponse struct {
	File models.File `json:"file"`
}

func (h FileHandler) UploadFileHandler(c *gin.Context) {
	cfg := config.LoadConfig()

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	destination := c.PostForm("destination")

	objectName := uuid.New().String()
	bucketName := cfg.MinioBucketName
	fileType := file.Header.Get("Content-Type")

	err = h.S3.UploadFile(file, bucketName, objectName, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	dbFile := models.File{
		Path:        "/" + objectName,
		Type:        fileType,
		Destination: destination,
	}

	err = h.DB.Create(&dbFile).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file information"})
		return
	}

	fullObjectPath := "http://" + cfg.MinioEndpoint + "/" + bucketName + "/" + objectName
	dbFile.Path = fullObjectPath

	c.JSON(http.StatusOK, FileUploadResponse{
		File: dbFile,
	})
}
