package database

import (
	config "learning-go/internal/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dbConfig *config.DBConfig) (db *gorm.DB, err error) {
	host := dbConfig.Host
	username := dbConfig.User
	password := dbConfig.Password
	dbname := dbConfig.Name
	port := string(dbConfig.Port)
	ssl := dbConfig.SSLMode
	timezone := dbConfig.Timezone

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
