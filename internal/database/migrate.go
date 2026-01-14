package database

import (
	order "github.com/QuangNV23062004/learning-go/internal/pkg/orders/domain"
	product "github.com/QuangNV23062004/learning-go/internal/pkg/products/domain"
	user "github.com/QuangNV23062004/learning-go/internal/pkg/users/domain"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&user.User{},
		&product.Product{},
		&order.Order{},
	)
	if err != nil {
		return err
	}
	return nil
}
