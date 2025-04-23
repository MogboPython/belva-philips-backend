package service

import (
	"encoding/json"
	"fmt"

	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type OrderService interface {
	CreateOrder(req *model.OrderRequest) (*model.OrderResponse, error)
	GetOrderByID(id string) (*model.OrderResponse, error)
	GetAllOrders(page, limit string) ([]*model.OrderResponse, error)
	GetOrdersByUserID(userID, pageStr, limitStr string) ([]*model.OrderResponse, error)
	UpdateOrderStatus(orderID string, request *model.OrderStatusChangeRequest) (*model.OrderResponse, error)
	// DeleteUser(id int64) error
}

type orderService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
}

func NewOrderService(orderRepo repository.OrderRepository, userRepo repository.UserRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

func (s *orderService) CreateOrder(request *model.OrderRequest) (*model.OrderResponse, error) {
	user, err := s.userRepo.GetByEmail(request.UserEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	detailsBytes, err := json.Marshal(request.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal details: %w", err)
	}

	detailsJSON := datatypes.JSON(detailsBytes)

	shotsArray := pq.StringArray(request.Shots)

	// Create a new order
	order := &model.Order{
		UserID:             user.ID,
		ProductName:        request.ProductName,
		ProductDescription: request.ProductDescription,
		Details:            detailsJSON,
		FinishType:         request.FinishType,
		Quantity:           request.Quantity,
		ShootType:          request.ShootType,
		Status:             request.Status,
		Shots:              shotsArray,
		DeliverySpeed:      request.DeliverySpeed,
	}

	// Save to database
	if err := s.orderRepo.Create(order); err != nil {
		log.Error("error saving user: %v", err)
		return nil, err
	}

	return mapOrderToResponse(order), nil
}

// GetAllUsers retrieves all orders
func (s *orderService) GetAllOrders(pageStr, limitStr string) ([]*model.OrderResponse, error) {
	// Convert to integers
	offset, limit := utils.GetPageAndLimitInt(pageStr, limitStr)

	orders, err := s.orderRepo.GetAll(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	orderResponses := make([]*model.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = mapOrderToResponse(order)
	}

	return orderResponses, nil
}

// GetOrderByID retrieves an order by ID
func (s *orderService) GetOrderByID(id string) (*model.OrderResponse, error) {
	order, err := s.orderRepo.GetByOrderID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	return mapOrderToResponse(order), nil
}

func (s *orderService) GetOrdersByUserID(userID, pageStr, limitStr string) ([]*model.OrderResponse, error) {
	offset, limit := utils.GetPageAndLimitInt(pageStr, limitStr)

	orders, err := s.orderRepo.GetByUserID(userID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find orders: %w", err)
	}

	orderResponses := make([]*model.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = mapOrderToResponse(order)
	}

	return orderResponses, nil
}

func (s *orderService) UpdateOrderStatus(orderID string, request *model.OrderStatusChangeRequest) (*model.OrderResponse, error) {
	order, err := s.orderRepo.UpdateOrder(orderID, request.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return mapOrderToResponse(order), nil
}

// mapOrderToResponse maps a order model to a order response
func mapOrderToResponse(order *model.Order) *model.OrderResponse {
	var detailsMap map[string]any

	shotsStringArray := []string(order.Shots)

	if err := json.Unmarshal(order.Details, &detailsMap); err != nil {
		detailsMap = nil

		log.Error("error unmarshaling order details: %v", err)
	}

	return &model.OrderResponse{
		ID:                   order.ID,
		UserID:               order.User.ID,
		UserEmail:            order.User.Email,
		UserMembershipStatus: order.User.MembershipStatus,
		ProductName:          order.ProductName,
		ProductDescription:   order.ProductDescription,
		ShootType:            order.ShootType,
		FinishType:           order.FinishType,
		Quantity:             order.Quantity,
		Details:              detailsMap,
		Shots:                shotsStringArray,
		DeliverySpeed:        order.DeliverySpeed,
		Status:               order.Status,
		CreatedAt:            order.CreatedAt,
		UpdatedAt:            order.UpdatedAt,

		// ProductDescriptionImage: order.ProductDescriptionImage,
	}
}
