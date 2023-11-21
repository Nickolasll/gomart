package domain

import "github.com/google/uuid"

type UserAggregateRepositoryInterface interface {
	Create(login string, password string) (UserAggregate, error)
	Get(userID uuid.UUID) (UserAggregate, error)
	GetByLogin(login string) (*UserAggregate, error)
	Save(userAggregate UserAggregate) error
	Init()
}

type BalanceRepositoryInterface interface {
	Get(userID uuid.UUID) (Balance, error)
}

type OrderRepositoryInterface interface {
	Get(number string) (*Order, error)
	GetAll(userID uuid.UUID) ([]Order, error)
}

type WithdrawRepositoryInterface interface {
	Get(number string) (*Withdraw, error)
	GetAll(userID uuid.UUID) ([]Withdraw, error)
}
