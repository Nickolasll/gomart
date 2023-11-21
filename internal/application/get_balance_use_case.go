package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
)

type GetBalance struct {
	balanceRepository domain.BalanceRepositoryInterface
}

func (u GetBalance) Execute(userID uuid.UUID) (domain.Balance, error) {
	balance, err := u.balanceRepository.Get(userID)
	if err != nil {
		return balance, err
	}
	return balance, nil
}
