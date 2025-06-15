package main

import (
	"log"

	"go-postgres-gorm-gin-api/bot"
	"go-postgres-gorm-gin-api/config"
	"go-postgres-gorm-gin-api/db"
	"go-postgres-gorm-gin-api/handlers"
	"go-postgres-gorm-gin-api/realtime"
	"go-postgres-gorm-gin-api/rtmp_proxy"
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

	realtime := realtime.NewRealtimeLinstance(postgres)
	realtime.Init()

	proxy := rtmp_proxy.NewProxy(":19350", "localhost:1935", postgres)
	err = proxy.Start()
	if err != nil {
		log.Fatal("failed to start rtmp proxy:", err)
	}
	defer proxy.Stop()

	adminApi := router.Group("/api")
	{
		tagHandler := handlers.NewTagHandler(postgres)
		adminApi.GET("/tags", tagHandler.GetTags)
		adminApi.POST("/tags", tagHandler.CreateTag)
		adminApi.GET("/tags/:id", tagHandler.GetTag)
		adminApi.PUT("/tags/:id", tagHandler.UpdateTag)
		adminApi.DELETE("/tags/:id", tagHandler.DeleteTag)

		mediaHandler := handlers.NewMediaHandler(postgres)
		adminApi.GET("/media", mediaHandler.GetResults)
		adminApi.GET("/media/:id", mediaHandler.GetMedia)
		adminApi.POST("/media", mediaHandler.CreateMedia)
		adminApi.DELETE("/media/:id", mediaHandler.DeleteMedia)
		adminApi.POST("/toggle-media-visibility/:id", mediaHandler.ToggleMediaVisibility)

		fileHandler := handlers.NewFileHandler(postgres, s3instance)
		adminApi.POST("/files/upload", fileHandler.UploadFileHandler)

		streamHandler := handlers.NewStreamHandler(postgres)
		adminApi.GET("/stream-status", streamHandler.GetStreamStatus)
		adminApi.GET("/stream-key", streamHandler.GetStreamKeyByRoom)
	}

	streamAuthApi := router.Group("/srs-api")
	{
		streamHandler := handlers.NewStreamHandler(postgres)

		streamAuthApi.POST("/auth", streamHandler.AuthentificateStreamHook)
	}

	publicApi := router.Group("/public")
	{
		mediaHandler := handlers.NewMediaHandler(postgres)
		publicApi.GET("/media", mediaHandler.GetMedias)

		streamHandler := handlers.NewStreamHandler(postgres)
		publicApi.GET("/hls-stream", streamHandler.GetStreamProxy)

		chatHandler := handlers.NewChatHandler(postgres, botInstance)
		publicApi.GET("/ws/chat", chatHandler.ConnectToChatWebSocket)
	}

	router.Run("0.0.0.0:" + cfg.Port)
}
