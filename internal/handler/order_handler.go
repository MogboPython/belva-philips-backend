package handler

import (
	"errors"
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/service"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderHandler struct {
	orderService service.OrderService
	validator    *validator.Validator
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		validator:    validator.New(),
	}
}

// CreateOrder creates a new order by a user
//
//	@Summary		Create a new order
//	@Description	Create a new order with the provided information
//	@Tags			orders
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.OrderRequest	true	"Order information"
//	@Success		201		{object}	model.ResponseHTTP{data=model.OrderResponse}
//	@Failure		400		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/orders [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var payload model.OrderRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Invalid request",
			Data:    nil,
		})
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	order, err := h.orderService.CreateOrder(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "failed to find user") {
			return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
				Success: false,
				Message: "User does not exists",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully saved order",
		Data:    *order,
	})
}

// GetAllOrders is a function to get all order data from the database
//
//	@Summary		Get all orders (strictly for admin)
//	@Description	Fetch a paginated list of orders from the database
//	@Tags			orders
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int		false	"Page number (default is 1)"
//	@Param			limit	query		int		false	"Number of orders per page (default is 10)"
//	@Param			status	query		string	false	"order status (active or pending)"
//	@Success		200		{array}		model.ResponseHTTP{data=[]model.TotalOrderResponse}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/orders [get]
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")
	status := c.Query("status", "")

	orders, err := h.orderService.GetAllOrders(pageStr, limitStr, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully retrieved orders.",
		Data:    orders,
	})
}

// GetOrderByID is a function to get an order by ID
//
//	@Summary		Get order by ID
//	@Description	Get order by ID
//	@Tags			orders
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Order ID"
//	@Success		200	{object}	model.ResponseHTTP{data=model.OrderResponse}
//	@Failure		404	{object}	model.ResponseHTTP{}
//	@Failure		500	{object}	model.ResponseHTTP{}
//	@Router			/api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")

	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ResponseHTTP{
				Success: false,
				Message: "Order not found",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully found order",
		Data:    *order,
	})
}

// GetOrdersByUserID is a function to get orders by a single user
//
//	@Summary		Get order by User ID
//	@Description	Get order by User ID
//	@Tags			orders
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		string	false	"User ID of the user"
//	@Param			page	query		int		false	"Page number (default is 1)"
//	@Param			limit	query		int		false	"Number of orders per page (default is 10)"
//	@Success		200		{object}	model.ResponseHTTP{data=model.OrderResponse}
//	@Failure		404		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/orders/user/{userId} [get]
func (h *OrderHandler) GetOrdersByUserID(c *fiber.Ctx) error {
	userID := c.Params("userId")
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	orders, err := h.orderService.GetOrdersByUserID(userID, pageStr, limitStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully retrieved orders.",
		Data:    orders,
	})
}

// UpdateOrderStatus is a function to update an order status
//
//	@Summary		Update the status of an order (strictly for admin)
//	@Description	Update the status of an order
//	@Tags			orders
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"Order ID"
//	@Param			status	body		model.OrderStatusChangeRequest	true	"Status update"
//	@Success		200		{object}	model.ResponseHTTP{data=model.OrderResponse}
//	@Failure		404		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/orders/{order_id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("order_id")

	var payload model.OrderStatusChangeRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	order, err := h.orderService.UpdateOrderStatus(id, &payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ResponseHTTP{
				Success: false,
				Message: "Order not found",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully updated order status",
		Data:    *order,
	})
}
