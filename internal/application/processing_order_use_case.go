package application

import (
	"errors"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
)

type ProcessingOrder struct {
	userAggregateRepository domain.UserAggregateRepositoryInterface
	accrualClient           domain.AccrualClientInterface
}

func (p ProcessingOrder) updateOrder(UserID uuid.UUID, order domain.Order) error {
	user, err := p.userAggregateRepository.Get(UserID)
	if err != nil {
		return err
	}
	user = user.UpdateOrder(order)
	p.userAggregateRepository.Save(user)
	return nil
}

func (p ProcessingOrder) Execute(order domain.Order) (bool, error) {
	AccrualOrderResponse, err := p.accrualClient.GetOrderStatus(order.Number)
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			order.Status = domain.StatusInvalid
			err = p.updateOrder(order.UserAggregateID, order)
			if err != nil {
				return false, err
			}
			return true, nil
		} else {
			return false, err
		}
	}
	order.Status = AccrualOrderResponse.Status
	if AccrualOrderResponse.Status == domain.StatusProcessed {
		order = order.SetAccrual(AccrualOrderResponse.Accrual)
	}
	err = p.updateOrder(order.UserAggregateID, order)
	if err != nil {
		return false, err
	}
	if order.Status == domain.StatusProcessing || order.Status == domain.StatusRegistred {
		return false, nil
	}
	return true, nil
}
