package infrastructure

import (
	"fmt"
	"learning-go/internal/infrastructure"
	"learning-go/internal/pkg/orders/domain"
	orderType "learning-go/internal/pkg/orders/types"
	"learning-go/internal/types"
	"math"

	"gorm.io/gorm"
)

type OrderRepository struct {
	*infrastructure.BaseRepository[*domain.Order]
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		BaseRepository: infrastructure.NewBaseRepository[*domain.Order](db),
		db:             db,
	}
}

func (r *OrderRepository) FindOrdersByUserIDWithOptions(id string, includeDeleted bool, options types.OrderOptions) ([]*orderType.OrderResponse, error) {
	var result []*orderType.OrderResponse

	where := r.db.Model(&domain.Order{}).Where("orders.user_id = ?", id)

	if !includeDeleted {
		where = where.Where("orders.is_deleted = ?", false)
	}

	selectFields := `
        orders.*`

	if options.WithUser {
		selectFields += `, users.username, users.email`
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}
	if options.WithProduct {
		selectFields += `, products.name AS product_name, products.price AS product_price, products.stock AS product_stock`
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	where = where.Select(selectFields)

	err := where.First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	return result, nil
}

func (r *OrderRepository) FindOrdersByUserIDPaginatedWithOptions(userID string, page int, limit int, search string, searchField string, order string, sortBy string, includeDeleted bool, options types.OrderOptions) (*types.Paginated[*orderType.OrderResponse], error) {
	var orders []*orderType.OrderResponse

	where := r.db.Model(&domain.Order{}).Where("orders.user_id = ?", userID)

	if !includeDeleted {
		where = where.Where("orders.is_deleted = ?", false)
	}

	selectFields := `
        orders.*`

	if options.WithUser {
		selectFields += `, users.username, users.email`
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}
	if options.WithProduct {
		selectFields += `, products.name AS product_name, products.price AS product_price, products.stock AS product_stock`
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	where = where.Select(selectFields)

	if search != "" && searchField != "" {
		where = where.Where(searchField+" ILIKE ?", "%"+search+"%")
	}

	// Count before pagination
	var total int64
	countQuery := r.db.Model(&domain.Order{}).Where("user_id = ?", userID)
	if !includeDeleted {
		countQuery = countQuery.Where("is_deleted = ?", false)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	safeSkip := (page - 1) * limit

	err := where.Order(sortBy + " " + order).Offset(safeSkip).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &types.Paginated[*orderType.OrderResponse]{
		Data:        orders,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		Order:       order,
		SortBy:      sortBy,
		HasPrevious: page > 1,
		HasNext:     page < totalPages,
	}, nil
}

func (r *OrderRepository) FindByIDWithOptions(id string, includeDeleted bool, options types.OrderOptions) (*orderType.OrderResponse, error) {
	var order *orderType.OrderResponse
	where := r.db.Model(&domain.Order{}).Where("orders.id = ?", id)
	if !includeDeleted {
		where = where.Where("orders.is_deleted = ?", false)
	}

	selectFields := `
        orders.*`

	if options.WithUser {
		selectFields += `, users.username, users.email`
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}
	if options.WithProduct {
		selectFields += `, products.name AS product_name, products.price AS product_price, products.stock AS product_stock`
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	where = where.Select(selectFields)

	err := where.First(&order).Error
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) PaginatedWithOptions(page int, limit int, search string, searchField string, order string, sortBy string, includeDeleted bool, options types.OrderOptions) (*types.Paginated[*orderType.OrderResponse], error) {
	var orders []*orderType.OrderResponse

	where := r.db.Model(new(domain.Order))

	if !includeDeleted {
		where = where.Where("orders.is_deleted = ?", false)
	}

	selectFields := `
        orders.*`

	if options.WithUser {
		selectFields += `, users.username, users.email`
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}
	if options.WithProduct {
		selectFields += `, products.name AS product_name, products.price AS product_price, products.stock AS product_stock`
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	where = where.Select(selectFields)

	if search != "" && searchField != "" {
		where = where.Where(fmt.Sprintf("%s ILIKE ?", searchField), "%"+search+"%")
	}

	// Count before pagination
	var total int64
	countQuery := r.db.Model(new(domain.Order))
	if !includeDeleted {
		countQuery = countQuery.Where("is_deleted = ?", false)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	skip := (page - 1) * limit

	if err := where.Order(sortBy + " " + order).Offset(skip).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &types.Paginated[*orderType.OrderResponse]{
		Data:        orders,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		Order:       order,
		SortBy:      sortBy,
		HasPrevious: page > 1,
		HasNext:     page < totalPages,
	}, nil
}

func (r *OrderRepository) FindAllWithOptions(includeDeleted bool, options types.OrderOptions) ([]*orderType.OrderResponse, error) {
	var orders []*orderType.OrderResponse
	where := r.db.Model(new(domain.Order))
	if !includeDeleted {
		where = where.Where("orders.is_deleted = ?", false)
	}

	selectFields := `
        orders.*`

	if options.WithUser {
		selectFields += `, users.username, users.email`
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}
	if options.WithProduct {
		selectFields += `, products.name AS product_name, products.price AS product_price, products.stock AS product_stock`
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	where = where.Select(selectFields)

	if err := where.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
