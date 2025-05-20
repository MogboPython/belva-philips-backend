package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	GetByOrderID(orderID string) (*model.Order, error)
	GetByUserID(userID string, offset, limit int) ([]*model.Order, error)
	Update(orderID, status string) (*model.Order, error)
	GetAll(offset, limit int, status string) ([]*model.Order, model.OrdersCount, error)
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
	exists, err := utils.ExistsByID(r.db, &model.User{}, order.UserID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("failed to find user")
	}

	orderName, err := r.generateUniqueOrderName()
	if err != nil {
		return fmt.Errorf("failed to generate order name: %w", err)
	}

	order.OrderName = orderName

	err = r.db.Create(&order).Error
	if err != nil {
		return err
	}

	return r.db.Model(&order).Association("User").Find(&order.User)
}

func (*orderRepository) generateUniqueOrderName() (string, error) {
	now := time.Now()
	date := now.Format("20060102")

	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("BELVA-%s-%s", date, id.String()[:6]), nil
}

func (r *orderRepository) GetByOrderID(orderID string) (*model.Order, error) {
	var order model.Order

	if err := r.db.Preload("User").Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetAll(offset, limit int, status string) ([]*model.Order, model.OrdersCount, error) {
	var orders []*model.Order

	var count model.OrdersCount

	tx := r.db.Model(&model.Order{})

	const (
		statusActive    = "quote_received"
		statusCompleted = "mark_completed"
	)

	switch status {
	case "active":
		tx = tx.Where("status = ?", statusActive)
	case "completed":
		tx = tx.Where("status = ?", statusCompleted)
	case "pending":
		tx = tx.Where("status != ? AND status != ?", statusActive, statusCompleted)
	}

	if err := tx.Preload("User").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, count, err
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Order{}).Count(&count.Total).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Order{}).Where("status = ?", statusActive).Count(&count.ActiveCount).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Order{}).Where("status = ?", statusCompleted).Count(&count.CompletedCount).Error; err != nil {
			return err
		}

		count.PendingCount = count.Total - count.ActiveCount - count.CompletedCount
		return nil
	}); err != nil {
		return nil, count, err
	}

	return orders, count, nil
}

func (r *orderRepository) GetByUserID(userID string, offset, limit int) ([]*model.Order, error) {
	var orders []*model.Order

	if err := r.db.Preload("User").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) Update(orderID, status string) (*model.Order, error) {
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
