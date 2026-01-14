package infrastructure

import (
	"github.com/QuangNV23062004/learning-go/internal/infrastructure"
	"github.com/QuangNV23062004/learning-go/internal/pkg/users/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	*infrastructure.BaseRepository[*domain.User]
	db *gorm.DB
}

// Constructor
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: infrastructure.NewBaseRepository[*domain.User](db),
		db:             db,
	}
}

func (r *UserRepository) FindByEmail(email string, tx *gorm.DB) (*domain.User, error) {
	user := &domain.User{}
	if err := r.GetDatabase(tx).Model(&domain.User{}).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	if user.IsDeleted {
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}
