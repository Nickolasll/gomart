package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserAggregate struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Login    string
	Password string
	Orders   []Order `gorm:"foreignKey:UserAggregateID"`
	Balance  Balance `gorm:"foreignKey:UserAggregateID"`
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
