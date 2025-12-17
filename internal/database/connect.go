package database

import (
	config "learning-go/internal/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (db *gorm.DB, err error) {
	host := config.GetEnv("DB_HOST", "localhost")
	username := config.GetEnv("DB_USER", "postgres")
	password := config.GetEnv("DB_PASSWORD", "postgres")
	dbname := config.GetEnv("DB_NAME", "postgres")
	port := config.GetEnv("DB_PORT", "5432")
	ssl := config.GetEnv("DB_SSLMODE", "disable")
	timezone := config.GetEnv("DB_TIMEZONE", "Asia/Ho_Chi_Minh")

	dns := "host=" + host + " user=" + username + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + ssl + " TimeZone=" + timezone
	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database !\n", err.Error())
	}

	log.Println("Connected to database successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Database logger set to Info level")

	return db, nil
}
