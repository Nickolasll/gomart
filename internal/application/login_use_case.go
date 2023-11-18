package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
)

type Login struct {
	userAggregateRepository domain.UserAggregateRepositoryInterface
	jose                    JOSEService
}

func (u Login) Execute(login string, password string) (string, error) {
	user, err := u.userAggregateRepository.GetByLogin(login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrLoginOrPasswordIsInvalid
	}
	if !u.jose.VerifyPassword(user.Password, password) {
		return "", ErrLoginOrPasswordIsInvalid
	}
	tokenString, err := u.jose.IssueToken(user.ID)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
