package database

import (
	"log"
	"sync"
	"time"

	"github.com/MrBista/blog-api/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func Connect() {
	once.Do(func() {
		dsn := config.AppConfig.DB.Dsn()

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("failed to get database object: %v", err)
		}

		sqlDB, err := db.DB()

		if err != nil {
			log.Fatalf("failed to get DB pool %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		DB = db

		log.Println("Successfullly to connect Database mysql")
	})
}

func Close() {
	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()

	if err != nil {
		log.Printf("failed to close connection db %v \n", err)
		return
	}

	sqlDB.Close()

}
