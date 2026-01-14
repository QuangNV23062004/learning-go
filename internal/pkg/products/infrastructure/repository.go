package infrastructure

import (
	"math"

	"github.com/QuangNV23062004/learning-go/internal/infrastructure"
	"github.com/QuangNV23062004/learning-go/internal/pkg/products/domain"
	"github.com/QuangNV23062004/learning-go/internal/types"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
	*infrastructure.BaseRepository[*domain.Product]
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db:             db,
		BaseRepository: infrastructure.NewBaseRepository[*domain.Product](db),
	}
}

func (r *ProductRepository) FindAllByUserID(userID string, includeDeleted bool, tx *gorm.DB) ([]*domain.Product, error) {
	var products []*domain.Product

	var where *gorm.DB = r.GetDatabase(tx)
	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}
	error := where.Model(&domain.Product{}).Where("user_id = ?", userID).Find(&products).Error
	if error != nil {
		return nil, error
	}

	return products, nil
}

func (r *ProductRepository) PaginatedByUserId(userID string, page int, limit int, search, searchField, order, sortBy string, includeDeleted bool, tx *gorm.DB) (*types.Paginated[domain.Product], error) {
	var products []domain.Product

	where := r.GetDatabase(tx).Model(&domain.Product{}).Where("user_id = ?", userID)
	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}

	if search != "" && searchField != "" {
		where = where.Where(searchField+" LIKE ?", "%"+search+"%")
	}

	safeSkip := (page - 1) * limit

	error := where.Order(sortBy + " " + order).Offset(safeSkip).Limit(limit).Find(&products).Error
	if error != nil {
		return nil, error
	}

	var total int64
	if err := where.Count(&total).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &types.Paginated[domain.Product]{
		Data:        products,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		Order:       order,
		SortBy:      sortBy,
		HasPrevious: page > 1,
		HasNext:     page < totalPages,
	}, nil

}
