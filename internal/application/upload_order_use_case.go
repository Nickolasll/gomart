package application

import (
	"sync"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
)

type UploadOrder struct {
	orderRepository         domain.OrderRepositoryInterface
	userAggregateRepository domain.UserAggregateRepositoryInterface
	ch                      chan<- domain.Order
	wg                      *sync.WaitGroup
}

func (u UploadOrder) Execute(userID uuid.UUID, number string) error {
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
	u.wg.Add(1)
	u.ch <- newOrder
	return nil
}
