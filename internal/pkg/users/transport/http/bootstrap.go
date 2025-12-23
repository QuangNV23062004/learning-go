package http

import (
	"learning-go/internal/config"
	"learning-go/internal/pkg/users/application"
	"learning-go/internal/pkg/users/infrastructure"
	"learning-go/internal/pkg/users/templates"
	"learning-go/internal/utils"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func BootstrapUserRoutes(api *fiber.App, db *gorm.DB, jwtConfig *config.JWTConfig, mailConfig *config.MailConfig, serverConfig *config.ServerConfig) {
	jwtService := utils.NewJwtService(jwtConfig)
	emailService := utils.NewEmailService(mailConfig, templates.TemplatesFS)
	passwordService := utils.NewPasswordService()
	userRepository := infrastructure.NewUserRepository(db)
	userService := application.NewUserService(userRepository, jwtService, emailService, passwordService, serverConfig)
	userHandler := NewUserHandler(userService)
	userRouter := NewRouter(userHandler, jwtService)
	userRouter.SetupRoutes(api)
}
