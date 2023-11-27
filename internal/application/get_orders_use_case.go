package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type getOrders struct {
	orderRepository domain.OrderRepositoryInterface
	log             *logrus.Logger
}

func (u getOrders) Execute(userID uuid.UUID) ([]domain.Order, error) {
	orders, err := u.orderRepository.GetAll(userID)
	if err != nil {
		return orders, err
	}
	for _, order := range orders {
		u.log.Info("+++")
		u.log.Info(order.UserAggregateID)
		u.log.Info(order.Number)
		u.log.Info(order.Status)
		u.log.Info(order.Accrual)
		u.log.Info("---")
	}
	return orders, nil
}
