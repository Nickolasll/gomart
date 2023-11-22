package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
)

type getBalance struct {
	balanceRepository domain.BalanceRepositoryInterface
}

func (u getBalance) Execute(userID uuid.UUID) (domain.Balance, error) {
	balance, err := u.balanceRepository.Get(userID)
	if err != nil {
		return balance, err
	}
	return balance, nil
}
