package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type uploadWithdraw struct {
	withdrawRepository      domain.WithdrawRepositoryInterface
	userAggregateRepository domain.UserAggregateRepositoryInterface
	log                     *logrus.Logger
}

func (u uploadWithdraw) Execute(userID uuid.UUID, number string, sum float64) error {
	if !IsValidNumber(number) {
		return ErrNotValidNumber
	}
	withdraw, err := u.withdrawRepository.Get(number)
	if err != nil {
		return err
	}

	if withdraw != nil {
		u.log.Info("Withdraw already uploaded")
		return ErrUploadedByAnotherUser
	}

	user, err := u.userAggregateRepository.Get(userID)
	if err != nil {
		return err
	}
	user, err = user.AddWithdraw(number, sum)
	if err != nil {
		return err
	}
	err = u.userAggregateRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}
