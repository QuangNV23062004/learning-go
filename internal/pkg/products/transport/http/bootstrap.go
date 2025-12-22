package http

import (
	"learning-go/internal/config"
	"learning-go/internal/pkg/products/application"
	"learning-go/internal/pkg/products/infrastructure"
	userInfrastructure "learning-go/internal/pkg/users/infrastructure"
	"learning-go/internal/utils"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func BootstrapProductRoutes(api *fiber.App, db *gorm.DB, jwtConfig *config.JWTConfig) {

	jwtService := utils.NewJwtService(jwtConfig)
	productRepository := infrastructure.NewProductRepository(db)
	userRepository := userInfrastructure.NewUserRepository(db)
	productService := application.NewProductService(productRepository, userRepository)
	productHandler := NewProductHandler(productService)
	productRouter := NewRouter(productHandler, jwtService)
	productRouter.SetupRoutes(api)
}
