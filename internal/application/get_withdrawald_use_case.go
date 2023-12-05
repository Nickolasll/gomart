package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
)

type getWithdrawals struct {
	withdrawRepository domain.WithdrawRepositoryInterface
}

func (u getWithdrawals) Execute(userID uuid.UUID) ([]domain.Withdraw, error) {
	withdrawals, err := u.withdrawRepository.GetAll(userID)
	if err != nil {
		return withdrawals, err
	}
	return withdrawals, nil
}
