package domain

import (
	"time"

	"github.com/google/uuid"
)

type Balance struct {
	UserAggregateID uuid.UUID `gorm:"type:uuid"`
	Current         int
	Withdraw        int
}

func (b Balance) CurrentToString() string {
	return MonetaryToString(b.Current)
}

func (b Balance) WithdrawToString() string {
	return MonetaryToString(b.Withdraw)
}

type Order struct {
	Number          string    `gorm:"primaryKey"`
	UserAggregateID uuid.UUID `gorm:"type:uuid"`
	Status          string
	UploadedAt      time.Time
	Accrual         int
}

func (o Order) AccrualToString() string {
	return MonetaryToString(o.Accrual)
}

func (o Order) SetAccrual(value string) Order {
	intValue := StringToMonetary(value)
	o.Accrual = intValue
	return o
}

type Withdraw struct {
	Order           string    `gorm:"primaryKey"`
	UserAggregateID uuid.UUID `gorm:"type:uuid"`
	Sum             int
	ProcessedAt     time.Time
}
