package infrastructure

import (
	"learning-go/internal/pkg/users/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// Constructor liked
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByID(id string, includeDeleted bool) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.Model(&domain.User{}).Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	if user.IsDeleted && !includeDeleted {
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (r *UserRepository) Create(user *domain.User) (*domain.User, error) {
	if err := r.db.Model(&domain.User{}).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *domain.User) (*domain.User, error) {
	user, err := r.FindByID(user.ID, false)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&domain.User{}).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Delete(id string) (bool, error) {
	user, err := r.FindByID(id, false)
	if err != nil {
		return false, err
	}

	user.IsDeleted = true
	user.DeletedAt = "now()"

	if err := r.db.Model(&domain.User{}).Save(user).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) Restore(id string) (bool, error) {

	user, err := r.FindByID(id, true)
	if err != nil {
		return false, err
	}

	user.IsDeleted = false
	user.DeletedAt = "null"

	if err := r.db.Model(&domain.User{}).Save(user).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.Model(&domain.User{}).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	if user.IsDeleted {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (r *UserRepository) FindAll(includeDeleted bool) ([]domain.User, error) {
	var users []domain.User
	var where *gorm.DB = r.db.Model(&domain.User{})
	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}

	if err := where.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) HardDelete(id string) (bool, error) {
	user, err := r.FindByID(id, true)
	if err != nil {
		return false, err
	}

	if err := r.db.Model(&domain.User{}).Unscoped().Delete(user).Error; err != nil {
		return false, err
	}

	return true, nil
}

// parameters should be validated in service layer
func (r *UserRepository) Paginated(page int, limit int, search string, searchField string, order string, sortBy string, includeDeleted bool) ([]domain.User, error) {

	var users []domain.User
	var where *gorm.DB = r.db.Model(&domain.User{})

	if search != "" {
		where = where.Where(searchField+" ILIKE ?", "%"+search+"%")
	}

	if !includeDeleted {
		where = where.Where("is_deleted = ?", false)
	}

	var take = limit
	var skip = (page - 1) * limit
	if error := where.Order(sortBy + " " + order).Take(take).Offset(skip).Find(&users).Error; error != nil {
		return nil, error
	}
	return users, nil
}
