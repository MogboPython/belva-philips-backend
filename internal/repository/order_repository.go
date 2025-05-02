package repository

import (
	"errors"
	"time"

	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	GetByOrderID(orderID string) (*model.Order, error)
	GetByUserID(userID string, offset, limit int) ([]*model.Order, error)
	UpdateOrder(orderID, status string) (*model.Order, error)
	GetAll(offset, limit int, status string) ([]*model.Order, OrdersCount, error)
	// Delete(id int64) error
}

type orderRepository struct {
	db *gorm.DB
}

type OrdersCount struct {
	Total        int64
	ActiveCount  int64
	PendingCount int64
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(order *model.Order) error {
	exists, err := utils.ExistsByID(r.db, &model.User{}, order.UserID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("failed to find user")
	}

	if err := r.db.Create(&order).Error; err != nil {
		return err
	}

	return r.db.Model(&order).Association("User").Find(&order.User)
}

func (r *orderRepository) GetByOrderID(orderID string) (*model.Order, error) {
	var order model.Order

	if err := r.db.Preload("User").Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// GetAll retrieves all orders
func (r *orderRepository) GetAll(offset, limit int, status string) ([]*model.Order, OrdersCount, error) {
	var orders []*model.Order

	var count OrdersCount

	var query *gorm.DB

	const (
		statusActive = "QUOTE RECEIVED"
	)

	// Building the query based on the status
	switch status {
	case "active":
		query = r.db.Preload("User").Where("status = ?", statusActive)
	case "pending":
		query = r.db.Preload("User").Where("status != ?", statusActive)
	default:
		query = r.db.Preload("User")
	}

	// Execute the query with pagination
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, count, err
	}

	// Get all counts in one database transaction
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get total count
		if err := tx.Model(&model.Order{}).Count(&count.Total).Error; err != nil {
			return err
		}

		// Get active count
		if err := tx.Model(&model.Order{}).Where("status = ?", statusActive).Count(&count.ActiveCount).Error; err != nil {
			return err
		}

		// Calculate pending count
		count.PendingCount = count.Total - count.ActiveCount
		return nil
	}); err != nil {
		return nil, count, err
	}

	return orders, count, nil
}

// GetByUserID retrieves orders by user ID
func (r *orderRepository) GetByUserID(userID string, offset, limit int) ([]*model.Order, error) {
	var orders []*model.Order

	if err := r.db.Preload("User").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) UpdateOrder(orderID, status string) (*model.Order, error) {
	var order model.Order

	if err := r.db.Preload("User").Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}

	order.Status = status
	order.UpdatedAt = time.Now()

	if err := r.db.Save(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}
