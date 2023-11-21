package infrastructure

import (
	"errors"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAggregateRepository struct {
	DB gorm.DB
}

func (u UserAggregateRepository) Init() {
	u.DB.AutoMigrate(
		&domain.UserAggregate{},
		&domain.Order{},
		&domain.Balance{},
		&domain.Withdraw{},
	)
}

func (u UserAggregateRepository) Create(login string, password string) (domain.UserAggregate, error) {
	userID := uuid.New()
	user := domain.UserAggregate{
		ID:       userID,
		Login:    login,
		Password: password,
		Balance:  domain.Balance{UserAggregateID: userID, Current: 0, Withdraw: 0},
	}
	err := u.DB.Create(&user).Error
	if err != nil {
		return user, err
	}
	err = u.DB.Save(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u UserAggregateRepository) Get(ID uuid.UUID) (domain.UserAggregate, error) {
	var user domain.UserAggregate
	err := u.DB.Preload("Balance").Preload("Orders").Preload("Withdrawals").First(&user, ID).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u UserAggregateRepository) GetByLogin(login string) (*domain.UserAggregate, error) {
	var user domain.UserAggregate
	err := u.DB.Where("login = ?", login).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			// Постоянно прокидывать ошибку выше или просто логгировать ее здесь?
			return nil, err
		}
	}
	return &user, nil
}

func (u UserAggregateRepository) Save(userAggregate domain.UserAggregate) error {
	err := u.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&userAggregate).Error
	if err != nil {
		return err
	}
	return nil
}
