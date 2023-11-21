package infrastructure

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BalanceRepository struct {
	DB gorm.DB
}

func (b BalanceRepository) Get(userID uuid.UUID) (domain.Balance, error) {
	var balance domain.Balance
	err := b.DB.First(&balance, userID).Error
	if err != nil {
		return balance, err
	}
	return balance, nil
}
