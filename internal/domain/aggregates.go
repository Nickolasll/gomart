package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInsufficientFunds = errors.New("insufficient funds on current user balance")

type UserAggregate struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Login       string
	Password    string
	Orders      []Order    `gorm:"foreignKey:UserAggregateID"`
	Balance     Balance    `gorm:"foreignKey:UserAggregateID;references:ID"`
	Withdrawals []Withdraw `gorm:"foreignKey:UserAggregateID"`
}

func (u UserAggregate) AddOrder(number string) (UserAggregate, Order) {
	order := Order{
		UserAggregateID: u.ID,
		Number:          number,
		UploadedAt:      time.Now(),
		Status:          StatusNew,
		Accrual:         0,
	}
	u.Orders = append(u.Orders, order)
	return u, order
}

func (u UserAggregate) UpdateOrder(order Order) UserAggregate {
	for i, o := range u.Orders {
		if o.Number == order.Number {
			u.Orders[i] = order
			if order.Status == StatusProcessed {
				u.Balance.Current += order.Accrual
			}
			return u
		}
	}
	return u
}

func (u UserAggregate) AddWithdraw(number string, sum float64) (UserAggregate, error) {
	withdraw := Withdraw{
		Order:           number,
		UserAggregateID: u.ID,
		ProcessedAt:     time.Now(),
	}
	withdraw = withdraw.SetSum(sum)
	// if u.Balance.Current < withdraw.Sum {
	// 	return u, ErrInsufficientFunds
	// }
	u.Withdrawals = append(u.Withdrawals, withdraw)
	u.Balance.Current -= withdraw.Sum
	u.Balance.Withdraw += withdraw.Sum
	return u, nil
}
