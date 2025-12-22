package main

import (
	"learning-go/internal/config"
	"learning-go/internal/database"
	orderTransport "learning-go/internal/pkg/orders/transport/http"
	productTransport "learning-go/internal/pkg/products/transport/http"
	userTransport "learning-go/internal/pkg/users/transport/http"

	"learning-go/internal/types"
	"log"
	"strconv"

	httpError "learning-go/internal/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/joho/godotenv"
)

func main() {
	//load env variables
	godotenv.Load()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			status := httpError.GetStatusCode(err)
			return c.Status(status).JSON(types.Response{
				Status:  status,
				Success: false,
				Error:   err.Error(),
			})
		}})

	app.Use(helmet.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", config.GetEnv("FRONTEND_URL", "http://localhost:3000")},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	app.Use(compress.New())

	//health check route
	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	appConfig := config.AppConfig{

		DB_CONFIG: config.DBConfig{
			Host:     config.GetEnv("DB_HOST", "localhost"),
			Port:     config.GetEnv("DB_PORT", "5432"),
			User:     config.GetEnv("DB_USER", "postgres"),
			Password: config.GetEnv("DB_PASSWORD", "postgres"),
			Name:     config.GetEnv("DB_NAME", "learning-go"),
			SSLMode:  config.GetEnv("DB_SSLMODE", "disable"),
			Timezone: config.GetEnv("DB_TIMEZONE", "Asia/Ho_Chi_Minh"),
		},

		JWT_CONFIG: config.JWTConfig{
			Issuer:        config.GetEnv("JWT_ISSUER", "my-golang-app"),
			AccessSecret:  config.GetEnv("JWT_ACCESS_SECRET", "acd7aa40a1f9c7df48adfaf48c69e4c3203a6a17acc60638a7bdc9103d0a499997180cb8d8"),
			RefreshSecret: config.GetEnv("JWT_REFRESH_SECRET", "8ee21eb52bf23900b62daa6fc7f7b8a09e6548cc8ef5ecc894fa26ff49022b8b8327ba6ce9914b"),
			AccessExpiry:  config.GetEnv("JWT_ACCESS_EXPIRY", "1h"),
			RefreshExpiry: config.GetEnv("JWT_REFRESH_EXPIRY", "1h"),
			VerifySecret:  config.GetEnv("JWT_VERIFY_SECRET", "5f4dcc3b5aa765d61d8327deb882cf99acd7aa40a1f9c7df48adfaf48c69e4c3203a6a17acc60638a7bdc9103d0a499997180cb8d8"),
			VerifyExpiry:  config.GetEnv("JWT_VERIFY_EXPIRY", "30m"),
		},

		MAIL_CONFIG: config.MailConfig{
			Host:     config.GetEnv("MAIL_HOST", "smtp.gmail.com"),
			Port:     config.GetEnvAsInt("MAIL_PORT", 587),
			Username: config.GetEnv("MAIL_USERNAME", "example@gmail.com"),
			Password: config.GetEnv("MAIL_PASSWORD", "xxxx xxxx xxxx xxxx"),
		},

		SERVER_CONFIG: config.ServerConfig{
			Host:    config.GetEnv("SERVER_HOST", "http://localhost:2000"),
			Port:    config.GetEnvAsInt("PORT", 2000),
			AppName: config.GetEnv("APP_NAME", "My Golang App"),
		},
	}

	//database connections
	db, err := database.Connect(&appConfig.DB_CONFIG)

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	//database migrations
	database.Migrate(db)

	//setup routes
	userTransport.BootstrapUserRoutes(app, db, &appConfig.JWT_CONFIG, &appConfig.MAIL_CONFIG, &appConfig.SERVER_CONFIG)
	productTransport.BootstrapProductRoutes(app, db, &appConfig.JWT_CONFIG)
	orderTransport.BootstrapOrderRoutes(app, db, &appConfig.JWT_CONFIG)

	//start server
	port := appConfig.SERVER_CONFIG.Port

	log.Printf("Server starting on :%d", port)
	if err := app.Listen(":" + strconv.Itoa(port)); err != nil {
		log.Fatal(err)
	}
}
