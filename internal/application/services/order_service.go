package services

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// OrderService implements the order service interface
type OrderService struct {
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo ports.OrderRepository, productRepo ports.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, customerID string, shippingInfo entities.ShippingInfo) (*entities.Order, error) {
	if customerID == "" {
		return nil, errors.New("customer ID is required")
	}

	order := entities.NewOrder(customerID, shippingInfo)
	err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, id string) (*entities.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

// AddItem adds a product variation to an order
func (s *OrderService) AddItem(ctx context.Context, orderID, productID, variationID string, quantity int) error {
	// Get the order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Get the product
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// Find the variation
	variation, err := product.GetVariation(variationID)
	if err != nil {
		return err
	}

	// Check stock availability
	if variation.StockQuantity < quantity {
		return errors.New("insufficient stock")
	}

	// Add item to order
	err = order.AddItem(productID, variationID, variation.SKU, product.Name, variation.Price, quantity)
	if err != nil {
		return err
	}

	// Update order in database
	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return err
	}

	// Update product stock
	variation.StockQuantity -= quantity
	return s.productRepo.UpdateStock(ctx, variationID, variation.StockQuantity)
}

// UpdateItemQuantity updates the quantity of an item in an order
func (s *OrderService) UpdateItemQuantity(ctx context.Context, orderID, productID, variationID string, quantity int) error {
	// Get the order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Find the current quantity of the item in the order
	currentQuantity := 0
	for _, item := range order.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			currentQuantity = item.Quantity
			break
		}
	}

	// If removing the item
	if quantity <= 0 {
		err = order.RemoveItem(productID, variationID)
		if err != nil {
			return err
		}

		// Update order in database
		err = s.orderRepo.Update(ctx, order)
		if err != nil {
			return err
		}

		// Update product stock
		product, err := s.productRepo.GetByID(ctx, productID)
		if err != nil {
			return err
		}

		variation, err := product.GetVariation(variationID)
		if err != nil {
			return err
		}

		variation.StockQuantity += currentQuantity
		return s.productRepo.UpdateStock(ctx, variationID, variation.StockQuantity)
	}

	// If increasing quantity, check stock availability
	if quantity > currentQuantity {
		product, err := s.productRepo.GetByID(ctx, productID)
		if err != nil {
			return err
		}

		variation, err := product.GetVariation(variationID)
		if err != nil {
			return err
		}

		additionalQuantity := quantity - currentQuantity
		if variation.StockQuantity < additionalQuantity {
			return errors.New("insufficient stock")
		}

		// Update product stock
		variation.StockQuantity -= additionalQuantity
		err = s.productRepo.UpdateStock(ctx, variationID, variation.StockQuantity)
		if err != nil {
			return err
		}
	} else if quantity < currentQuantity {
		// If decreasing quantity, return stock
		product, err := s.productRepo.GetByID(ctx, productID)
		if err != nil {
			return err
		}

		variation, err := product.GetVariation(variationID)
		if err != nil {
			return err
		}

		returnedQuantity := currentQuantity - quantity
		variation.StockQuantity += returnedQuantity
		err = s.productRepo.UpdateStock(ctx, variationID, variation.StockQuantity)
		if err != nil {
			return err
		}
	}

	// Update item quantity in order
	err = order.UpdateItemQuantity(productID, variationID, quantity)
	if err != nil {
		return err
	}

	// Update order in database
	return s.orderRepo.Update(ctx, order)
}

// RemoveItem removes an item from an order
func (s *OrderService) RemoveItem(ctx context.Context, orderID, productID, variationID string) error {
	return s.UpdateItemQuantity(ctx, orderID, productID, variationID, 0)
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, status entities.OrderStatus) error {
	// Get the order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Update status
	err = order.UpdateStatus(status)
	if err != nil {
		return err
	}

	// Update order in database
	return s.orderRepo.Update(ctx, order)
}

// SetPaymentInfo sets the payment information for an order
func (s *OrderService) SetPaymentInfo(ctx context.Context, orderID, method, transactionID string, amount float64) error {
	// Get the order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Set payment info
	err = order.SetPaymentInfo(method, transactionID, amount)
	if err != nil {
		return err
	}

	// Update order in database
	return s.orderRepo.Update(ctx, order)
}

// SetTrackingInfo sets the shipping tracking information for an order
func (s *OrderService) SetTrackingInfo(ctx context.Context, orderID, carrier, trackingNum string) error {
	// Get the order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Set tracking info
	err = order.SetTrackingInfo(carrier, trackingNum)
	if err != nil {
		return err
	}

	// Update order in database
	return s.orderRepo.Update(ctx, order)
}

// ListOrders retrieves orders with filtering and pagination
func (s *OrderService) ListOrders(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Order, int, error) {
	return s.orderRepo.List(ctx, filter, page, limit)
}

// GetCustomerOrders retrieves orders for a specific customer
func (s *OrderService) GetCustomerOrders(ctx context.Context, customerID string, page, limit int) ([]*entities.Order, int, error) {
	return s.orderRepo.GetByCustomerID(ctx, customerID, page, limit)
}
