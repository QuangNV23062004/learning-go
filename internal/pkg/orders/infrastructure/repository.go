package infrastructure

import (
	"fmt"
	"learning-go/internal/infrastructure"
	"learning-go/internal/pkg/orders/domain"
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
		db: db,
	}
}

func (r *OrderRepository) FindOrdersByUserIDWithOptions(userID string, includeDeleted bool, options types.OrderOptions) ([]*domain.Order, error) {
	var orders []*domain.Order

	var where *gorm.DB = r.db.Model(&domain.Order{})

	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}

	if options.WithUser == true {
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}

	if options.WithProduct == true {
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	err := where.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) FindOrdersByUserIDPaginatedWithOptions(userID string, page int, limit int, search string, searchField string, order string, sortBy string, includeDeleted bool, options types.OrderOptions) (*types.Paginated[*domain.Order], error) {
	var orders []*domain.Order

	where := r.db.Model(&domain.Order{}).Where("user_id = ?", userID)

	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}
	if search != "" && searchField != "" {
		where = where.Where(searchField+" ILIKE ?", "%"+search+"%")
	}

	if options.WithUser == true {
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}

	if options.WithProduct == true {
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	safeSkip := (page - 1) * limit

	err := where.Order(sortBy + " " + order).Offset(safeSkip).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	count := where.Model(&domain.Order{}).Where("user_id = ?", userID)
	if !includeDeleted {
		count = count.Where("is_deleted = ?", false)
	}
	var total int64
	err = count.Count(&total).Error
	if err != nil {
		return nil, err
	}

	var totalPages int = int(math.Ceil((float64(total)) / float64(limit)))

	return &types.Paginated[*domain.Order]{
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

func (r *OrderRepository) FindByIDWithOptions(id string, includeDeleted bool, options types.OrderOptions) (*domain.Order, error) {
	var order *domain.Order
	where := r.db.Model(&domain.Order{})
	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}

	if options.WithUser == true {
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}

	if options.WithProduct == true {
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	err := where.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) PaginatedWithOptions(page int, limit int, search string, searchField string, order string, sortBy string, includeDeleted bool, options types.OrderOptions) (*types.Paginated[*domain.Order], error) {
	var entities []*domain.Order
	where := r.db.Model(new(domain.Order))

	if search != "" && searchField != "" {
		where = where.Where(fmt.Sprintf("%s ILIKE ?", searchField), "%"+search+"%")
	}
	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}

	take := limit
	skip := (page - 1) * limit

	if options.WithUser == true {
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}

	if options.WithProduct == true {
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	if err := where.Order(sortBy + " " + order).Offset(skip).Limit(take).Find(&entities).Error; err != nil {
		return nil, err
	}

	var total int64
	if err := where.Count(&total).Error; err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &types.Paginated[*domain.Order]{
		Data:        entities,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		Order:       order,
		SortBy:      sortBy,
		HasPrevious: page > 1,
		HasNext:     page < totalPages,
	}, nil
}

func (r *OrderRepository) FindAllWithOptions(includeDeleted bool, options types.OrderOptions) ([]*domain.Order, error) {
	var entities []*domain.Order
	where := r.db.Model(new(domain.Order))
	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}

	if options.WithUser == true {
		where = where.Joins("LEFT JOIN users ON users.id = orders.user_id")
	}

	if options.WithProduct == true {
		where = where.Joins("LEFT JOIN products ON products.id = orders.product_id")
	}

	if err := where.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}
