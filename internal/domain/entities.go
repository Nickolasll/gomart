package domain

import (
	"time"

	"github.com/google/uuid"
)

// По поводу полей типа time.Time.
// Просто я думал на будущее, вдруг захотим фильтровать заказы и списания по дате
// Я понимаю, что мне так проще, но - это менее оптимально
// Если будет нужно, я переделаю

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
	if u.Balance.Current < withdraw.Sum {
		return u, ErrInsufficientFunds
	}
	u.Withdrawals = append(u.Withdrawals, withdraw)
	u.Balance.Current -= withdraw.Sum
	u.Balance.Withdrawn += withdraw.Sum
	return u, nil
}

type Balance struct {
	UserAggregateID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Current         int
	Withdrawn       int
}

func (b Balance) CurrentToFloat() float64 {
	return MonetaryToFloat(b.Current)
}

func (b Balance) WithdrawnToFloat() float64 {
	return MonetaryToFloat(b.Withdrawn)
}

type Order struct {
	Number          string    `gorm:"primaryKey"`
	UserAggregateID uuid.UUID `gorm:"type:uuid"`
	Status          string
	UploadedAt      time.Time
	Accrual         int
}

func (o Order) AccrualToFloat() float64 {
	return MonetaryToFloat(o.Accrual)
}

func (o Order) SetAccrual(value float64) Order {
	intValue := FloatToMonetary(value)
	o.Accrual = intValue
	return o
}

type Withdraw struct {
	Order           string    `gorm:"primaryKey"`
	UserAggregateID uuid.UUID `gorm:"type:uuid"`
	Sum             int
	ProcessedAt     time.Time
}

func (w Withdraw) SetSum(value float64) Withdraw {
	intValue := FloatToMonetary(value)
	w.Sum = intValue
	return w
}

func (w Withdraw) SumToFloat() float64 {
	return MonetaryToFloat(w.Sum)
}
