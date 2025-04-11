package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// UserService interface defines methods for user business logic
type OrderService interface {
	CreateOrder(req *model.OrderRequest) (*model.OrderResponse, error)
	GetOrderByID(id string) (*model.OrderResponse, error)
	GetAllOrders(page, limit string) ([]*model.OrderResponse, error)
	GetOrdersByUserID(user_id, pageStr, limitStr string) ([]*model.OrderResponse, error)
	// UpdateOrder(id int64, req *model.UpdateUserRequest) (*model.UserResponse, error)
	// DeleteUser(id int64) error
}

// userService implements UserService interface
type orderService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
}

// NewUserService creates a new user service
func NewOrderService(orderRepo repository.OrderRepository, userRepo repository.UserRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

func (s *orderService) CreateOrder(request *model.OrderRequest) (*model.OrderResponse, error) {
	// Validate common required fields
	if request.UserEmail == "" || request.ProductName == "" || request.ShootType == "" {
		return nil, errors.New("missing required fields: UserEmail, ProductName, or ShootType")
	}

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
		log.Printf("error saving user: %v", err)
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

	var orderResponses []*model.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, mapOrderToResponse(order))
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

func (s *orderService) GetOrdersByUserID(user_id, pageStr, limitStr string) ([]*model.OrderResponse, error) {
	offset, limit := utils.GetPageAndLimitInt(pageStr, limitStr)

	orders, err := s.orderRepo.GetByUserID(user_id, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find orders: %w", err)
	}

	var orderResponses []*model.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, mapOrderToResponse(order))
	}

	return orderResponses, nil
}

// mapOrderToResponse maps a order model to a order response
func mapOrderToResponse(order *model.Order) *model.OrderResponse {
	var detailsMap map[string]interface{}
	if err := json.Unmarshal(order.Details, &detailsMap); err != nil {
		log.Printf("error unmarshaling order details: %v", err)
		detailsMap = nil
	}
	shotsStringArray := []string(order.Shots)
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
