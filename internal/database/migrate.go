package database

import (
	order "learning-go/internal/pkg/orders/domain"
	product "learning-go/internal/pkg/products/domain"
	user "learning-go/internal/pkg/users/domain"

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
