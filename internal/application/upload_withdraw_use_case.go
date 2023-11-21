package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UploadWithdraw struct {
	withdrawRepository      domain.WithdrawRepositoryInterface
	userAggregateRepository domain.UserAggregateRepositoryInterface
	log                     *logrus.Logger
}

func (u UploadWithdraw) Execute(userID uuid.UUID, number string, sum float64) error {
	if !IsValidNumber(number) {
		u.log.Info("ErrNotValidNumber")
		return ErrNotValidNumber
	}
	withdraw, err := u.withdrawRepository.Get(number)
	if err != nil {
		u.log.Info("get withdraw " + err.Error())
		return err
	}

	if withdraw != nil {
		u.log.Info("withdraw exists")
		return ErrUploadedByAnotherUser
	}

	user, err := u.userAggregateRepository.Get(userID)
	if err != nil {
		u.log.Info("Get user err " + err.Error())
		return err
	}
	u.log.Info(user.Balance.Current)
	u.log.Info(user.Balance.Withdraw)
	user, err = user.AddWithdraw(number, sum)
	if err != nil {
		u.log.Info("add wighdraw " + err.Error())
		return err
	}
	u.log.Info(user.Balance.Current)
	u.log.Info(user.Balance.Withdraw)
	err = u.userAggregateRepository.Save(user)
	if err != nil {
		u.log.Info("save " + err.Error())
		return err
	}
	return nil
}
