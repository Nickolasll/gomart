package infrastructure

import (
	"errors"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB gorm.DB
}

func (o OrderRepository) Get(number string) (*domain.Order, error) {
	var order domain.Order
	err := o.DB.First(&order, number).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			// Постоянно прокидывать ошибку выше или просто логгировать ее здесь?
			return nil, err
		}
	}
	return &order, nil
}

func (o OrderRepository) GetAll(userID uuid.UUID) ([]domain.Order, error) {
	var orders []domain.Order
	err := o.DB.Where("user_aggregate_id = ?", userID).Find(&orders).Error
	if err != nil {
		return orders, err
	}
	return orders, nil
}
