package repository

import (
	"time"

	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	GetByOrderID(orderID string) (*model.Order, error)
	GetByUserID(userID string, offset, limit int) ([]*model.Order, error)
	UpdateOrder(orderID, status string) (*model.Order, error)
	GetAll(page, limit int) ([]*model.Order, error)
	// Delete(id int64) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(order *model.Order) error {
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
func (r *orderRepository) GetAll(offset, limit int) ([]*model.Order, error) {
	var orders []*model.Order

	if err := r.db.Preload("User").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

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
