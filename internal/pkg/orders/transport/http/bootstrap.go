package http

import (
	"github.com/QuangNV23062004/learning-go/internal/config"
	"github.com/QuangNV23062004/learning-go/internal/pkg/orders/application"
	"github.com/QuangNV23062004/learning-go/internal/pkg/orders/infrastructure"
	productInfrastructure "github.com/QuangNV23062004/learning-go/internal/pkg/products/infrastructure"
	userInfrastructure "github.com/QuangNV23062004/learning-go/internal/pkg/users/infrastructure"
	"github.com/QuangNV23062004/learning-go/internal/utils"

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
