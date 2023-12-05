package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/sirupsen/logrus"
)

type registration struct {
	userAggregateRepository domain.UserAggregateRepositoryInterface
	jose                    JOSEService
	log                     *logrus.Logger
}

func (u registration) Execute(login string, password string) (string, error) {
	user, err := u.userAggregateRepository.GetByLogin(login)
	if err != nil {
		return "", err
	}
	if user != nil {
		u.log.Info("User already exists")
		return "", ErrLoginAlreadyInUse
	}
	hashedPassword := u.jose.Hash(password)
	userAggregate, err := u.userAggregateRepository.Create(login, hashedPassword)
	if err != nil {
		return "", err
	}
	tokenString, err := u.jose.IssueToken(userAggregate.ID)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
