package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/111zxc/cocktaildle/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Не удалось настроить пул соединений: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = DB.AutoMigrate(&models.User{}, &models.DailyGame{}, &models.UserStats{}, &models.GameAttempt{}, &models.Guess{})
	if err != nil {
		log.Fatalf("Не удалось выполнить миграции: %v", err)
	}

	log.Println("Подключение к базе данных успешно установлено и схема инициализирована")
}
