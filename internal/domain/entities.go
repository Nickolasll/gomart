package domain

import (
	"time"

	"github.com/google/uuid"
)

type Balance struct {
	UserAggregateID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Current         int
	Withdraw        int
}

func (b Balance) CurrentToFloat() float32 {
	return MonetaryToFloat(b.Current)
}

func (b Balance) WithdrawToFloat() float32 {
	return MonetaryToFloat(b.Withdraw)
}

type Order struct {
	Number          string    `gorm:"primaryKey"`
	UserAggregateID uuid.UUID `gorm:"type:uuid"`
	Status          string
	UploadedAt      time.Time
	Accrual         int
}

func (o Order) AccrualToFloat() float32 {
	return MonetaryToFloat(o.Accrual)
}

func (o Order) SetAccrual(value float32) Order {
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

func (w Withdraw) SetSum(value float32) Withdraw {
	intValue := FloatToMonetary(value)
	w.Sum = intValue
	return w
}

func (w Withdraw) SumToFloat() float32 {
	return MonetaryToFloat(w.Sum)
}
