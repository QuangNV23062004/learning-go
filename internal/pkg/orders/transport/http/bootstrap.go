package http

import (
	"learning-go/internal/config"
	"learning-go/internal/pkg/orders/application"
	"learning-go/internal/pkg/orders/infrastructure"
	productInfrastructure "learning-go/internal/pkg/products/infrastructure"
	userInfrastructure "learning-go/internal/pkg/users/infrastructure"
	"learning-go/internal/utils"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func BootstrapOrderRoutes(api *fiber.App, db *gorm.DB, jwtConfig *config.JWTConfig) {

	jwtService := utils.NewJwtService(jwtConfig)
	repo := infrastructure.NewOrderRepository(db)
	userRepository := userInfrastructure.NewUserRepository(db)
	productRepository := productInfrastructure.NewProductRepository(db)
	orderService := application.NewOrderService(repo, userRepository, productRepository)
	orderHandler := NewOrderHandler(orderService)
	orderRouter := NewRouter(orderHandler, jwtService)
	orderRouter.SetupRoutes(api)
}
