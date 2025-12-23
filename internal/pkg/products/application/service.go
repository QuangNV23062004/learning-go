package application

import (
	"learning-go/internal/http"
	"learning-go/internal/pkg/products/domain"
	"learning-go/internal/pkg/products/dtos"
	"learning-go/internal/pkg/products/infrastructure"
	roleEnums "learning-go/internal/pkg/users/enums"
	userInfrastructure "learning-go/internal/pkg/users/infrastructure"
	"learning-go/internal/types"

	"gorm.io/gorm"
)

type ProductService struct {
	repo     *infrastructure.ProductRepository
	userRepo *userInfrastructure.UserRepository
}

func NewProductService(repo *infrastructure.ProductRepository, userRepo *userInfrastructure.UserRepository) *ProductService {
	return &ProductService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *ProductService) GetProductByID(id string, includeDeleted bool, role string) (*domain.Product, error) {

	safeIncludeDeleted := false

	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	product, err := s.repo.FindByID(id, safeIncludeDeleted, nil)
	if product == nil {
		return nil, domain.ErrProductNotFound
	}

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetProductsByUserID(userID string, includeDeleted bool, role string) ([]*domain.Product, error) {

	users, err := s.userRepo.FindByID(userID, false, nil)
	if users == nil {
		return nil, domain.ErrUserNotFound
	}

	safeIncludeDeleted := false

	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	products, err := s.repo.FindAllByUserID(userID, safeIncludeDeleted, nil)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetPaginatedProductsByUserID(userID string, page int, limit int, search, searchField, order, sortBy string, includeDeleted bool, role string) (*types.Paginated[domain.Product], error) {
	users, err := s.userRepo.FindByID(userID, false, nil)
	if users == nil {
		return nil, domain.ErrUserNotFound
	}

	safeIncludeDeleted := false

	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	paginatedProducts, err := s.repo.PaginatedByUserId(userID, page, limit, search, searchField, order, sortBy, safeIncludeDeleted, nil)
	if err != nil {
		return nil, err
	}
	return paginatedProducts, nil
}

func (s *ProductService) GetAllProducts(includeDeleted bool, role string) ([]*domain.Product, error) {

	safeIncludeDeleted := false

	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	products, err := s.repo.FindAll(safeIncludeDeleted, nil)

	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetPaginatedProducts(page int, limit int, search, searchField, order, sortBy string, includeDeleted bool, role string) (*types.Paginated[*domain.Product], error) {
	safeIncludeDeleted := false

	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	paginatedProducts, err := s.repo.Paginated(page, limit, search, searchField, order, sortBy, safeIncludeDeleted, nil)
	if err != nil {
		return nil, err
	}
	return paginatedProducts, nil
}

func (s *ProductService) CreateProduct(product *dtos.CreateProductDTO) (*domain.Product, error) {
	db := s.repo.GetDatabase(nil)
	var createdProduct *domain.Product
	err := db.Transaction(func(tx *gorm.DB) error {
		user, err := s.userRepo.FindByID(product.UserID, false, tx)

		if user == nil {
			return domain.ErrUserNotFound
		}

		if err != nil {
			return err
		}

		newProduct := &domain.Product{
			Name:   product.Name,
			Price:  product.Price,
			Stock:  product.Stock,
			UserID: product.UserID,
		}

		createdProduct, err = s.repo.Create(newProduct, tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (s *ProductService) UpdateProduct(id string, product *dtos.UpdateProductDTO, UserID string, role string) (*domain.Product, error) {

	var updatedProduct *domain.Product

	db := s.repo.GetDatabase(nil)

	err := db.Transaction(func(tx *gorm.DB) error {

		existingProduct, err := s.repo.FindByID(id, false, tx)

		if existingProduct == nil {
			return domain.ErrProductNotFound
		}

		if err != nil {
			return err
		}

		if (UserID != existingProduct.UserID) && (role != string(roleEnums.Admin)) {
			return http.ErrForbidden
		}

		if product.Name != "" {
			existingProduct.Name = product.Name
		}

		if product.Price != 0 {
			existingProduct.Price = product.Price
		}

		if product.Stock != 0 {
			existingProduct.Stock = product.Stock
		}

		updatedProduct, err = s.repo.Update(existingProduct, tx)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return updatedProduct, nil
}

func (s *ProductService) DeleteProduct(id string, role string, UserID string) (bool, error) {
	var deleted bool
	db := s.repo.GetDatabase(nil)

	err := db.Transaction(func(tx *gorm.DB) error {
		existingProduct, err := s.repo.FindByID(id, false, tx)
		if existingProduct == nil {
			return domain.ErrProductNotFound
		}

		if err != nil {
			return err
		}

		if (UserID != existingProduct.UserID) && (role != string(roleEnums.Admin)) {
			return http.ErrForbidden
		}

		deleted, err = s.repo.Delete(id, tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return deleted, nil

}

func (s *ProductService) RestoreProduct(id string) (bool, error) {
	var restored bool
	db := s.repo.GetDatabase(nil)
	err := db.Transaction(func(tx *gorm.DB) error {
		deletedProduct, err := s.repo.FindByID(id, true, tx)
		if deletedProduct == nil {
			return domain.ErrProductNotFound
		}
		if err != nil {
			return err
		}
		restored, err = s.repo.Restore(id, tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return restored, nil
}
