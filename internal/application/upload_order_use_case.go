package application

import (
	"sync"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
)

type uploadOrder struct {
	orderRepository         domain.OrderRepositoryInterface
	userAggregateRepository domain.UserAggregateRepositoryInterface
	ch                      chan<- domain.Order
	wg                      *sync.WaitGroup
	w                       Worker
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
	u.wg.Add(1)
	go u.w.Serve()
	return nil
}
