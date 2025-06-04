package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config stores the application configuration
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	Port       string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	TelegramBotToken string
	TelegramChatId   int64

	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

var ConfigInstance Config

// LoadConfig loads the configuration from environment variables
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: ", err)
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		dbPortStr = "5432"
	}
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "postgres"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = ""
	}

	redisDB := 0
	redisDBStr := os.Getenv("REDIS_DB")
	if redisDBStr != "" {
		redisDB, err = strconv.Atoi(redisDBStr)
	}

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		telegramBotToken = ""
	}

	telegramChatId := os.Getenv("TELEGRAM_CHAT_ID")
	if telegramChatId == "" {
		telegramChatId = ""
	}

	telegramChatIdInt, err := strconv.ParseInt(telegramChatId, 10, 64)
	if err != nil {
		log.Fatalf("Invalid TELEGRAM_CHAT_ID: %v", err)
	}

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "127.0.0.1:9000"
	}

	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioBucketName := os.Getenv("MINIO_BUCKET_NAME")

	ConfigInstance = Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,

		Port: port,

		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,

		TelegramBotToken: telegramBotToken,
		TelegramChatId:   telegramChatIdInt,

		MinioEndpoint:   minioEndpoint,
		MinioAccessKey:  minioAccessKey,
		MinioSecretKey:  minioSecretKey,
		MinioBucketName: minioBucketName,
	}

	return ConfigInstance
}

func GetConfig() *Config {
	return &ConfigInstance
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}
