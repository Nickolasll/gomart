package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
)

type getOrders struct {
	orderRepository domain.OrderRepositoryInterface
}

func (u getOrders) Execute(userID uuid.UUID) ([]domain.Order, error) {
	orders, err := u.orderRepository.GetAll(userID)
	if err != nil {
		return orders, err
	}
	return orders, nil
}
