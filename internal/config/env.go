package config

import (
	"os"
	"strconv"
)

type AppConfig struct {
	DB_CONFIG DBConfig

	//jwt
	JWT_CONFIG JWTConfig

	//mail
	MAIL_CONFIG MailConfig
	//server
	SERVER_CONFIG ServerConfig
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SSLMode  string
	Timezone string
}

type JWTConfig struct {
	Issuer        string
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  string
	RefreshExpiry string
	VerifySecret  string
	VerifyExpiry  string
}

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type ServerConfig struct {
	Host    string
	Port    int
	AppName string
}

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultVal
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultVal
	}
	return value
}
