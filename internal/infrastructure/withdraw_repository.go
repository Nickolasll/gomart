package infrastructure

import (
	"errors"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WithdrawRepository struct {
	DB gorm.DB
}

func (w WithdrawRepository) Get(number string) (*domain.Withdraw, error) {
	var withdraw domain.Withdraw
	err := w.DB.First(&withdraw, number).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &withdraw, nil
}

func (w WithdrawRepository) GetAll(userID uuid.UUID) ([]domain.Withdraw, error) {
	var withdrawals []domain.Withdraw
	err := w.DB.Where("user_aggregate_id = ?", userID).Find(&withdrawals).Error
	if err != nil {
		return withdrawals, err
	}
	return withdrawals, nil
}
