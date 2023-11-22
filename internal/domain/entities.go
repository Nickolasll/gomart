package domain

import (
	"time"

	"github.com/google/uuid"
)

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
