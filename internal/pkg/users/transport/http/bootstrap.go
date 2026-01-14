package http

import (
	"github.com/QuangNV23062004/learning-go/internal/config"
	"github.com/QuangNV23062004/learning-go/internal/pkg/users/application"
	"github.com/QuangNV23062004/learning-go/internal/pkg/users/infrastructure"
	"github.com/QuangNV23062004/learning-go/internal/pkg/users/templates"
	"github.com/QuangNV23062004/learning-go/internal/utils"

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
