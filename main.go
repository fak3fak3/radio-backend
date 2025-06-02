package main

import (
	"log"

	"go-postgres-gorm-gin-api/bot"
	"go-postgres-gorm-gin-api/config"
	"go-postgres-gorm-gin-api/db"
	"go-postgres-gorm-gin-api/handlers"
	"go-postgres-gorm-gin-api/s3"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	postgres, err := db.ConnectAndMigratePostgres(&cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}

	_, err = db.ConnectRedis(&cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
		return
	}

	s3instance, err := s3.ConnectMinio(&cfg)
	if err != nil {
		log.Fatal("Failed to connect to Minio:", err)
		return
	}

	botInstance := bot.Init(&cfg)
	go botInstance.Start()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	chatHandler := handlers.NewChatHandler(postgres, botInstance)
	router.GET("/ws/chat", chatHandler.ConnectToChatWebSocket)

	apiGroup := router.Group("/api")
	{
		tagHandler := handlers.NewTagHandler(postgres)
		apiGroup.GET("/tags", tagHandler.GetTags)
		apiGroup.POST("/tags", tagHandler.CreateTag)
		apiGroup.GET("/tags/:id", tagHandler.GetTag)
		apiGroup.PUT("/tags/:id", tagHandler.UpdateTag)
		apiGroup.DELETE("/tags/:id", tagHandler.DeleteTag)

		mediaHandler := handlers.NewMediaHandler(postgres)
		apiGroup.GET("/media", mediaHandler.GetResults)
		apiGroup.POST("/media", mediaHandler.CreateMedia)

		fileHandler := handlers.NewFileHandler(postgres, s3instance)
		apiGroup.POST("/files/upload", fileHandler.UploadFileHandler)
	}

	router.Run(":" + cfg.Port)
}
