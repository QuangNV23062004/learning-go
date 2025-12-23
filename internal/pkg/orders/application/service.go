package application

import (
	"learning-go/internal/pkg/orders/domain"
	"learning-go/internal/pkg/orders/dtos"
	"learning-go/internal/pkg/orders/infrastructure"
	productInfrastructure "learning-go/internal/pkg/products/infrastructure"
	roleEnums "learning-go/internal/pkg/users/enums"
	userInfrastructure "learning-go/internal/pkg/users/infrastructure"
	"learning-go/internal/types"
	"log"

	orderType "learning-go/internal/pkg/orders/types"
	productDomain "learning-go/internal/pkg/products/domain"

	"gorm.io/gorm"
)

type OrderService struct {
	repo        *infrastructure.OrderRepository
	userRepo    *userInfrastructure.UserRepository
	productRepo *productInfrastructure.ProductRepository
}

func NewOrderService(repo *infrastructure.OrderRepository, userRepo *userInfrastructure.UserRepository, productRepo *productInfrastructure.ProductRepository) *OrderService {
	return &OrderService{
		repo:        repo,
		userRepo:    userRepo,
		productRepo: productRepo,
	}
}

// only admin and owner can see order
func (s *OrderService) FindOrderByID(id string, includeDeleted bool, role string, sub string) (*orderType.OrderResponse, error) {

	safeIncludeDeleted := false
	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	order, err := s.repo.FindByIDWithOptions(id, safeIncludeDeleted, types.OrderOptions{
		WithUser:    true,
		WithProduct: true,
	})

	if order == nil {
		return nil, domain.ErrOrderNotFound
	}

	if role != string(roleEnums.Admin) && order.UserID != sub {
		return nil, domain.ErrNotAllowed
	}

	if err != nil {
		return nil, err
	}

	return order, nil
}

// only admin and owner can see orders
func (s *OrderService) FindOrdersByUserID(userID string, includeDeleted bool, sub string, role string) ([]*orderType.OrderResponse, error) {

	if role != string(roleEnums.Admin) && sub != userID {
		return nil, domain.ErrNotAllowed
	}

	users, err := s.userRepo.FindByID(userID, false, nil)
	if users == nil {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	safeIncludeDeleted := false

	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	orders, err := s.repo.FindOrdersByUserIDWithOptions(userID, safeIncludeDeleted, types.OrderOptions{
		WithUser:    true,
		WithProduct: true,
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// only admin and owner can see paginated orders
func (s *OrderService) FindPaginatedOrdersByUserID(userID string, page int, limit int, search string, searchField string, order string, sortBy string, includeDeleted bool, sub string, role string) (*types.Paginated[*orderType.OrderResponse], error) {

	if role != string(roleEnums.Admin) && sub != userID {
		return nil, domain.ErrNotAllowed
	}

	users, err := s.userRepo.FindByID(userID, false, nil)
	if users == nil {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	safeIncludeDeleted := false

	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	paginatedOrders, err := s.repo.FindOrdersByUserIDPaginatedWithOptions(userID, page, limit, search, searchField, order, sortBy, safeIncludeDeleted, types.OrderOptions{
		WithUser:    true,
		WithProduct: true,
	})

	if err != nil {
		return nil, err
	}

	return paginatedOrders, nil
}

// only admin can see all orders
func (s *OrderService) FindAllOrders(includeDeleted bool, role string) ([]*orderType.OrderResponse, error) {
	safeIncludeDeleted := false
	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted
	}

	orders, err := s.repo.FindAllWithOptions(safeIncludeDeleted, types.OrderOptions{
		WithUser:    true,
		WithProduct: true,
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// only admin can see paginated all orders
func (s *OrderService) Paginated(page int, limit int, search string, searchField string, order string, sortBy string, includeDeleted bool, role string) (*types.Paginated[*orderType.OrderResponse], error) {
	safeIncludeDeleted := false
	if role == string(roleEnums.Admin) {
		safeIncludeDeleted = includeDeleted

	}

	paginatedOrders, err := s.repo.PaginatedWithOptions(page, limit, search, searchField, order, sortBy, safeIncludeDeleted, types.OrderOptions{
		WithUser:    true,
		WithProduct: true,
	})

	if err != nil {
		return nil, err
	}
	return paginatedOrders, nil
}

func (s *OrderService) Create(orderDto *dtos.CreateOrderDTO, sub string) (*domain.Order, error) {
	db := s.repo.GetDatabase(nil)

	var createdOrder *domain.Order

	err := db.Transaction(func(tx *gorm.DB) error {
		user, err := s.userRepo.FindByID(sub, false, tx)
		if user == nil {
			return domain.ErrUserNotFound
		}

		if err != nil {
			return err
		}

		product, err := s.productRepo.FindByID(orderDto.ProductID, false, tx)
		if product == nil {
			return domain.ErrProductNotFound
		}

		if err != nil {
			return err
		}

		if product.Stock < orderDto.Quantity {
			return domain.ErrInsufficientStock
		}

		orderData := &domain.Order{
			UserID:    sub,
			ProductID: orderDto.ProductID,
			Quantity:  orderDto.Quantity,
			Total:     float64(orderDto.Quantity) * product.Price,
		}

		createdOrder, err = s.repo.Create(orderData, tx)
		if err != nil {
			return err
		}

		// Decrease product stock
		product.Stock -= orderDto.Quantity
		_, err = s.productRepo.Update(product, tx)
		if err != nil {
			return err
		}
		log.Print(product)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

// The api will handle two update scenarios, but irl product change should be rare if not impossible
// case 1: change product => restock old product, decrease new product stock
// case 2: change quantity only => adjust stock accordingly
func (s *OrderService) Update(id string, orderDto *dtos.UpdateOrderDTO, sub string, role string) (*domain.Order, error) {

	var updated *domain.Order

	db := s.repo.GetDatabase(nil)

	err := db.Transaction(func(tx *gorm.DB) error {
		//check order
		order, err := s.repo.FindByID(id, false, tx)
		if order == nil {
			return domain.ErrOrderNotFound
		}

		if err != nil {
			return err
		}

		//flags
		changeProductFlag := false

		if order.ProductID != orderDto.ProductID {
			changeProductFlag = true
		}

		// only admin and owner can update order
		if role != string(roleEnums.Admin) && order.UserID != sub {
			return domain.ErrNotAllowed
		}

		//user check
		user, err := s.userRepo.FindByID(sub, false, tx)
		if user == nil {
			return domain.ErrUserNotFound
		}

		if err != nil {
			return err
		}

		//product check
		product, err := s.productRepo.FindByID(orderDto.ProductID, false, tx)
		if product == nil {
			return domain.ErrProductNotFound
		}

		if err != nil {
			return err
		}

		//check & restock old product
		var oldProduct *productDomain.Product = nil
		if changeProductFlag {
			oldProduct, err = s.productRepo.FindByID(order.ProductID, true, tx)
			if oldProduct == nil {
				return domain.ErrProductNotFound
			}

			if err != nil {
				return err
			}

		}

		if oldProduct != nil {
			oldProduct.Stock += order.Quantity
			_, err = s.productRepo.Update(oldProduct, tx)
			if err != nil {
				return err
			}
		}

		availableStock := product.Stock

		//if same product => account current order quantity, else ignore
		if !changeProductFlag {
			availableStock = product.Stock + order.Quantity
		}

		if availableStock < orderDto.Quantity {
			return domain.ErrInsufficientStock
		}

		product.Stock = availableStock - orderDto.Quantity
		_, err = s.productRepo.Update(product, tx)
		if err != nil {
			return err
		}

		if orderDto.ProductID != "" {
			order.ProductID = orderDto.ProductID
		}
		if orderDto.Quantity != 0 {
			order.Quantity = orderDto.Quantity
		}
		order.Total = float64(orderDto.Quantity) * product.Price
		updated, err = s.repo.Update(order, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return updated, nil

}

func (s *OrderService) Delete(id string, sub string, role string) (bool, error) {

	var deleted bool

	db := s.repo.GetDatabase(nil)

	err := db.Transaction(func(tx *gorm.DB) error {
		//check order
		order, err := s.repo.FindByID(id, false, tx)
		if order == nil {
			return domain.ErrOrderNotFound
		}

		if err != nil {
			return err
		}

		// only admin and owner can delete order
		if role != string(roleEnums.Admin) && order.UserID != sub {
			return domain.ErrNotAllowed
		}

		//restock product
		product, err := s.productRepo.FindByID(order.ProductID, true, tx)
		if err != nil {
			return err
		}

		product.Stock += order.Quantity
		_, err = s.productRepo.Update(product, tx)
		if err != nil {
			return err
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

func (s *OrderService) Restore(id string) (bool, error) {
	var restored bool

	db := s.repo.GetDatabase(nil)

	err := db.Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.FindByID(id, true, tx)
		if order != nil {
			return domain.ErrOrderNotFound
		}

		if err != nil {
			return err
		}

		product, err := s.productRepo.FindByID(order.ProductID, true, tx)
		if product == nil {
			return domain.ErrProductNotFound
		}

		if err != nil {
			return err
		}

		if product.Stock < order.Quantity {
			return domain.ErrInsufficientStock
		}

		product.Stock -= order.Quantity

		updatedProduct, err := s.productRepo.Update(product, tx)
		if updatedProduct == nil {
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
