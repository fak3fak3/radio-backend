package db

import (
	"fmt"
	"go-postgres-gorm-gin-api/config"
	"go-postgres-gorm-gin-api/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectAndMigratePostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.DBHost + " port=" + fmt.Sprintf("%d", cfg.DBPort) + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	log.Println("Connected to Postgres database")
	DB = db

	// AutoMigrate the schema
	db.AutoMigrate(
		&models.Media{},
		&models.Tag{},
		&models.File{},
		&models.StreamData{},
		&models.StreamCredentials{},
		&models.StreamSession{},
	)

	return db, nil
}

func GetPostgresDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database connection is not initialized")
	}
	return DB
}
