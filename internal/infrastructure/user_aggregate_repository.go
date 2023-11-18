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
	u.DB.AutoMigrate(&domain.UserAggregate{})
	u.DB.AutoMigrate(&domain.Order{})
}

func (u UserAggregateRepository) Create(login string, password string) (domain.UserAggregate, error) {
	var user = domain.UserAggregate{ID: uuid.New(), Login: login, Password: password}
	res := u.DB.Create(&user)
	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}

func (u UserAggregateRepository) Get(ID uuid.UUID) (domain.UserAggregate, error) {
	var user domain.UserAggregate
	err := u.DB.First(&user, ID).Error
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
