package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
)

type uploadOrder struct {
	orderRepository         domain.OrderRepositoryInterface
	userAggregateRepository domain.UserAggregateRepositoryInterface
	ch                      chan<- domain.Order
}

func (u uploadOrder) Execute(userID uuid.UUID, number string) error {
	if !IsValidNumber(number) {
		return ErrNotValidNumber
	}
	order, err := u.orderRepository.Get(number)
	if err != nil {
		return err
	}

	if order != nil {
		if order.UserAggregateID == userID {
			return ErrUploadedByThisUser
		} else {
			return ErrUploadedByAnotherUser
		}
	}

	user, err := u.userAggregateRepository.Get(userID)
	if err != nil {
		return err
	}
	user, newOrder := user.AddOrder(number)
	err = u.userAggregateRepository.Save(user)
	if err != nil {
		return err
	}
	u.ch <- newOrder
	return nil
}
