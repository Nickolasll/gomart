package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type GetBalance struct {
	balanceRepository domain.BalanceRepositoryInterface
	log               *logrus.Logger
}

func (u GetBalance) Execute(userID uuid.UUID) (domain.Balance, error) {
	balance, err := u.balanceRepository.Get(userID)
	u.log.Info("GetBalance use case")
	u.log.Info(balance.Current)
	u.log.Info(balance.Withdraw)
	u.log.Info(balance.UserAggregateID.String())

	if err != nil {
		return balance, err
	}
	return balance, nil
}
