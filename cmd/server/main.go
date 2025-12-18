package main

import (
	"learning-go/internal/config"
	"learning-go/internal/database"
	"learning-go/internal/pkg/users/transport/http"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	database.Migrate(db)

	http.BootstrapUserRoutes(app, db)
	port := config.GetEnv("PORT", "2000")

	log.Printf("Server starting on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
