package infrastructure

import (
	"fmt"
	"learning-go/internal/domain"
	"learning-go/internal/types"
	"math"
	"time"

	"gorm.io/gorm"
)

type BaseModel interface {
	GetBaseEntity() *domain.BaseEntity
}

type BaseRepository[T BaseModel] struct {
	db *gorm.DB
}

func NewBaseRepository[T BaseModel](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) FindByID(id string, includeDeleted bool) (T, error) {
	var entity T
	if err := r.db.Where("id = ?", id).First(&entity).Error; err != nil {
		var zero T
		return zero, err
	}
	if be := entity.GetBaseEntity(); be != nil && be.IsDeleted && !includeDeleted {
		var zero T
		return zero, gorm.ErrRecordNotFound
	}
	return entity, nil
}

func (r *BaseRepository[T]) Create(entity T) (T, error) {
	if err := r.db.Create(&entity).Error; err != nil {
		var zero T
		return zero, err
	}
	return entity, nil
}

func (r *BaseRepository[T]) Update(entity T) (T, error) {
	be := entity.GetBaseEntity()
	if be == nil {
		var zero T
		return zero, fmt.Errorf("missing base entity")
	}
	if _, err := r.FindByID(be.ID, false); err != nil {
		var zero T
		return zero, err
	}
	if err := r.db.Model(&entity).Where("id = ?", be.ID).Updates(&entity).Error; err != nil {
		var zero T
		return zero, err
	}
	return entity, nil
}

func (r *BaseRepository[T]) Delete(id string) (bool, error) {
	entity, err := r.FindByID(id, false)
	if err != nil {
		return false, err
	}
	if be := entity.GetBaseEntity(); be != nil {
		be.IsDeleted = true
		be.DeletedAt = time.Now().Format(time.RFC3339)
	}
	if err := r.db.Model(&entity).Where("id = ?", id).Save(&entity).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *BaseRepository[T]) Restore(id string) (bool, error) {
	entity, err := r.FindByID(id, true)
	if err != nil {
		return false, err
	}
	if be := entity.GetBaseEntity(); be != nil {
		be.IsDeleted = false
		be.DeletedAt = ""
	}
	if err := r.db.Model(&entity).Where("id = ?", id).Save(&entity).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *BaseRepository[T]) FindAll(includeDeleted bool) ([]T, error) {
	var entities []T
	where := r.db.Model(new(T))
	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}
	if err := where.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *BaseRepository[T]) HardDelete(id string) (bool, error) {
	entity, err := r.FindByID(id, true)
	if err != nil {
		return false, err
	}
	if err := r.db.Unscoped().Delete(&entity).Error; err != nil {
		return false, err
	}
	return true, nil
}

// parameters should be validated in service layer
func (r *BaseRepository[T]) Paginated(page, limit int, search, searchField, order, sortBy string, includeDeleted bool) (*types.Paginated[T], error) {
	var entities []T
	where := r.db.Model(new(T))

	if search != "" && searchField != "" {
		where = where.Where(fmt.Sprintf("%s ILIKE ?", searchField), "%"+search+"%")
	}
	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}

	take := limit
	skip := (page - 1) * limit

	if err := where.Order(sortBy + " " + order).Offset(skip).Limit(take).Find(&entities).Error; err != nil {
		return nil, err
	}

	var total int64
	if err := where.Count(&total).Error; err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &types.Paginated[T]{
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
