package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
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
	GetAllOrders(page, limit, status string) (model.TotalOrderResponse, error)
	GetOrdersByUserID(userID, pageStr, limitStr string) ([]*model.OrderResponse, error)
	UpdateOrderStatus(orderID string, request *model.OrderStatusChangeRequest) (*model.OrderResponse, error)
	// TODO: DeleteOrder(id int64) error
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

func (s *orderService) CreateOrder(request *model.OrderRequest) (*model.OrderResponse, error) {
	var mails = 2

	detailsBytes, err := json.Marshal(request.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal details: %w", err)
	}

	detailsJSON := datatypes.JSON(detailsBytes)

	shotsArray := pq.StringArray(request.Shots)

	order := &model.Order{
		UserID:             request.UserID,
		ProductName:        request.ProductName,
		ProductDescription: request.ProductDescription,
		Details:            detailsJSON,
		FinishType:         request.FinishType,
		Quantity:           request.Quantity,
		ShootType:          request.ShootType,
		Status:             request.Status,
		MembershipType:     request.MembershipType,
		Shots:              shotsArray,
		DeliverySpeed:      request.DeliverySpeed,
	}

	if err := s.orderRepo.Create(order); err != nil {
		log.Error("error saving order: ", err)
		return nil, err
	}

	s.sendOrderNotificationEmailsConcurrently(order, mails)

	return mapOrderToResponse(order), nil
}

func (*orderService) sendOrderNotificationEmailsConcurrently(order *model.Order, mailsToSend int) {
	// TODO: read this
	var wg sync.WaitGroup

	wg.Add(mailsToSend)

	go func() {
		defer wg.Done()

		to := order.User.Email
		subject := "Your Quote Request at BelvaPhilips Imagery - Confirmation!"
		data := map[string]string{
			"Name":        order.User.Name,
			"ProductName": order.ProductName,
			"OrderDate":   order.CreatedAt.Format(time.UnixDate),
		}
		body, err := utils.ParseTemplate("user_confirmation.html", data)

		if err != nil {
			log.Errorf("Failed to parse user confirmation email template for Order ID %v: %v", order.ID, err)
			return
		}

		_, err = utils.SendEmail(to, subject, body)
		if err != nil {
			log.Errorf("Failed to send confirmation email to user %s for Order ID %v: %v", to, order.ID, err)
			return
		}

		log.Infof("Successfully sent confirmation email to user %s for Order ID: %v", to, order.ID)
	}()

	go func() {
		defer wg.Done()

		toAdmin := config.Config("ADMIN_EMAIL")
		subjectAdmin := "New Quote Request - Immediate Action Required!"
		dataAdmin := map[string]string{
			"OrderID":     order.ID,
			"ProductName": order.ProductName,
			"OrderDate":   order.CreatedAt.Format(time.UnixDate),
		}
		bodyAdmin, err := utils.ParseTemplate("order_notification.html", dataAdmin)

		if err != nil {
			log.Errorf("Failed to parse admin notification email template for Order ID %v: %v", order.ID, err)
			return
		}

		_, err = utils.SendEmail(toAdmin, subjectAdmin, bodyAdmin)
		if err != nil {
			log.Errorf("Failed to send notification email to admin %s for Order ID %v: %v", toAdmin, order.ID, err)
			return
		}

		log.Infof("Successfully sent notification email to admin %s for Order ID: %v", toAdmin, order.ID)
	}()

	wg.Wait()
}

func (s *orderService) GetAllOrders(pageStr, limitStr, status string) (model.TotalOrderResponse, error) {
	var totalOrderResponse model.TotalOrderResponse

	offset, limit := utils.GetPageAndLimitInt(pageStr, limitStr)

	orders, ordersCount, err := s.orderRepo.GetAll(offset, limit, status)
	if err != nil {
		return totalOrderResponse, fmt.Errorf("failed to get orders: %w", err)
	}

	formattedOrderResponses := make([]*model.OrderResponse, len(orders))
	for i, order := range orders {
		formattedOrderResponses[i] = mapOrderToResponse(order)
	}

	totalOrderResponse.Orders = formattedOrderResponses
	totalOrderResponse.OrderNumbers = ordersCount

	return totalOrderResponse, nil
}

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
	order, err := s.orderRepo.Update(orderID, request.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return mapOrderToResponse(order), nil
}

func mapOrderToResponse(order *model.Order) *model.OrderResponse {
	var detailsMap map[string]any

	shotsStringArray := []string(order.Shots)

	if err := json.Unmarshal(order.Details, &detailsMap); err != nil {
		detailsMap = nil

		log.Error("error unmarshaling order details: %v", err)
	}

	return &model.OrderResponse{
		ID:                   order.ID,
		OrderName:            order.OrderName,
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
		MembershipType:       order.MembershipType,
		CreatedAt:            order.CreatedAt,
		UpdatedAt:            order.UpdatedAt,
	}
}
