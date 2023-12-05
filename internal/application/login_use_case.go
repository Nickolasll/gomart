package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/sirupsen/logrus"
)

type login struct {
	userAggregateRepository domain.UserAggregateRepositoryInterface
	jose                    JOSEService
	log                     *logrus.Logger
}

func (u login) Execute(login string, password string) (string, error) {
	user, err := u.userAggregateRepository.GetByLogin(login)
	if err != nil {
		return "", err
	}
	if user == nil {
		u.log.Info("User with login " + login + " not found.")
		return "", ErrLoginOrPasswordIsInvalid
	}
	if !u.jose.VerifyPassword(user.Password, password) {
		u.log.Info("Invalid password for user " + login)
		return "", ErrLoginOrPasswordIsInvalid
	}
	tokenString, err := u.jose.IssueToken(user.ID)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
