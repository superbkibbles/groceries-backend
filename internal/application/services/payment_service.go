package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService struct {
	methodRepo   ports.PaymentMethodRepository
	gatewayRepo  ports.PaymentGatewayRepository
	customerRepo ports.CustomerPaymentMethodRepository
	orderRepo    ports.OrderRepository
}

func NewPaymentService(
	methodRepo ports.PaymentMethodRepository,
	gatewayRepo ports.PaymentGatewayRepository,
	customerRepo ports.CustomerPaymentMethodRepository,
	orderRepo ports.OrderRepository,
) *PaymentService {
	return &PaymentService{
		methodRepo:   methodRepo,
		gatewayRepo:  gatewayRepo,
		customerRepo: customerRepo,
		orderRepo:    orderRepo,
	}
}

// --- Payment Methods ---

func (s *PaymentService) CreatePaymentMethod(ctx context.Context, name, description string, methodType entities.PaymentMethodType, config map[string]interface{}) (*entities.PaymentMethod, error) {
	method := entities.NewPaymentMethod(name, description, methodType, config)
	if err := s.methodRepo.Create(ctx, method); err != nil {
		return nil, err
	}
	return method, nil
}

func (s *PaymentService) GetPaymentMethod(ctx context.Context, id string) (*entities.PaymentMethod, error) {
	return s.methodRepo.GetByID(ctx, id)
}

func (s *PaymentService) UpdatePaymentMethod(ctx context.Context, method *entities.PaymentMethod) error {
	return s.methodRepo.Update(ctx, method)
}

func (s *PaymentService) DeletePaymentMethod(ctx context.Context, id string) error {
	return s.methodRepo.Delete(ctx, id)
}

func (s *PaymentService) ListPaymentMethods(ctx context.Context, active bool) ([]*entities.PaymentMethod, error) {
	return s.methodRepo.List(ctx, active)
}

// --- Payment Gateways ---

func (s *PaymentService) CreatePaymentGateway(ctx context.Context, name, provider string, config map[string]interface{}) (*entities.PaymentGateway, error) {
	gateway := entities.NewPaymentGateway(name, provider, config)
	if err := s.gatewayRepo.Create(ctx, gateway); err != nil {
		return nil, err
	}
	return gateway, nil
}

func (s *PaymentService) GetPaymentGateway(ctx context.Context, id string) (*entities.PaymentGateway, error) {
	return s.gatewayRepo.GetByID(ctx, id)
}

func (s *PaymentService) UpdatePaymentGateway(ctx context.Context, gateway *entities.PaymentGateway) error {
	return s.gatewayRepo.Update(ctx, gateway)
}

func (s *PaymentService) DeletePaymentGateway(ctx context.Context, id string) error {
	return s.gatewayRepo.Delete(ctx, id)
}

func (s *PaymentService) ListPaymentGateways(ctx context.Context, active bool) ([]*entities.PaymentGateway, error) {
	return s.gatewayRepo.List(ctx, active)
}

// --- Customer Payment Methods ---

func (s *PaymentService) CreateCustomerPaymentMethod(ctx context.Context, customerID, paymentMethodID, token string, last4 string, expiryMonth, expiryYear int, isDefault bool) (*entities.CustomerPaymentMethod, error) {
	if isDefault {
		_ = s.customerRepo.ClearDefault(ctx, customerID)
	}

	method := entities.NewCustomerPaymentMethod(customerID, paymentMethodID, token, last4, expiryMonth, expiryYear, isDefault)
	if err := s.customerRepo.Create(ctx, method); err != nil {
		return nil, err
	}
	return method, nil
}

func (s *PaymentService) GetCustomerPaymentMethod(ctx context.Context, id string) (*entities.CustomerPaymentMethod, error) {
	return s.customerRepo.GetByID(ctx, id)
}

func (s *PaymentService) UpdateCustomerPaymentMethod(ctx context.Context, method *entities.CustomerPaymentMethod) error {
	return s.customerRepo.Update(ctx, method)
}

func (s *PaymentService) DeleteCustomerPaymentMethod(ctx context.Context, id string) error {
	return s.customerRepo.Delete(ctx, id)
}

func (s *PaymentService) ListCustomerPaymentMethods(ctx context.Context, customerID string) ([]*entities.CustomerPaymentMethod, error) {
	return s.customerRepo.GetByCustomer(ctx, customerID)
}

func (s *PaymentService) SetDefaultCustomerPaymentMethod(ctx context.Context, customerID, methodID string) error {
	if err := s.customerRepo.ClearDefault(ctx, customerID); err != nil {
		return err
	}
	return s.customerRepo.SetDefault(ctx, methodID)
}

// --- Payment Processing ---

func (s *PaymentService) ProcessPayment(ctx context.Context, orderID string, paymentMethodID string, amount float64) (*entities.PaymentInfo, error) {
	orderObjectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.orderRepo.GetByID(ctx, orderObjectID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.Status != entities.OrderStatusPending {
		return nil, errors.New("order is not in a payable state")
	}

	_, err = s.methodRepo.GetByID(ctx, paymentMethodID)
	if err != nil {
		return nil, errors.New("payment method not found")
	}

	transactionID := fmt.Sprintf("txn_%s_%d", orderID, time.Now().UnixNano())

	paymentInfo := entities.PaymentInfo{
		Method:        paymentMethodID,
		TransactionID: transactionID,
		PaidAt:        time.Now(),
		Status:        "paid",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	order.PaymentInfo = paymentInfo
	if err := order.UpdateStatus(entities.OrderStatusPaid); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return &paymentInfo, nil
}

func (s *PaymentService) VerifyPayment(ctx context.Context, transactionID string) (bool, error) {
	if transactionID == "" {
		return false, errors.New("transaction ID is required")
	}
	// In a real implementation this would query the payment gateway.
	// For now, we verify that the transaction ID looks valid.
	return len(transactionID) > 0, nil
}

func (s *PaymentService) RefundPayment(ctx context.Context, orderID string, amount float64, reason string) error {
	orderObjectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	order, err := s.orderRepo.GetByID(ctx, orderObjectID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.PaymentInfo.Status != "paid" {
		return errors.New("payment has not been completed")
	}

	order.PaymentInfo.Status = "refunded"
	order.PaymentInfo.Timestamp = time.Now()
	order.UpdatedAt = time.Now()

	return s.orderRepo.Update(ctx, order)
}
