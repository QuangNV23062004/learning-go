package http

import (
	"learning-go/internal/pkg/users/application"
	"learning-go/internal/pkg/users/infrastructure"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func BootstrapUserRoutes(api *fiber.App, db *gorm.DB) {
	userRepository := infrastructure.NewUserRepository(db)
	userService := application.NewUserService(userRepository)
	userHandler := NewUserHandler(userService)
	userRouter := NewRouter(userHandler)
	userRouter.SetupRoutes(api)
}
