package utils

import "go-postgres-gorm-gin-api/config"

func GenerateFullFilePath(objectName string) string {
	cfg := config.GetConfig()
	return "http://" + cfg.MinioEndpoint + "/" + cfg.MinioBucketName + objectName
}
