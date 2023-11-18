package infrastructure

import (
	"errors"

	"github.com/Nickolasll/gomart/internal/domain"
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
