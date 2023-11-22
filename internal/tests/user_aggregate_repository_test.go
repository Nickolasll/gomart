package tests

import (
	"testing"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestReadWrite(t *testing.T) {
	login := "login"
	password := "password"
	number := "12345"
	sum := float64(500.0)

	sqlDB, err := sql.Open("pgx", "postgresql://admin:admin@localhost:5432/postgres")
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(t, err)
	repo := infrastructure.UserAggregateRepository{DB: *db}
	repo.Init()
	err = db.Delete(&domain.Withdraw{Order: number}).Error
	require.NoError(t, err)
	err = db.Delete(&domain.Order{Number: number}).Error
	require.NoError(t, err)

	user, err := repo.Create(login, password)
	require.NoError(t, err)
	user, order := user.AddOrder(number)
	order.Status = domain.StatusProcessed
	order = order.SetAccrual(sum)
	user = user.UpdateOrder(order)
	user, err = user.AddWithdraw(number, sum)
	require.NoError(t, err)
	err = repo.Save(user)
	require.NoError(t, err)

	loadedUser, err := repo.Get(user.ID)
	require.NoError(t, err)
	assert.Equal(t, loadedUser.Login, login)
	assert.Equal(t, loadedUser.Password, password)
	assert.Equal(t, loadedUser.Balance.Current, 0)
	assert.Equal(t, loadedUser.Balance.Withdrawn, domain.FloatToMonetary(sum))
	loadedOrder := loadedUser.Orders[0]
	originalOrder := user.Orders[0]
	assert.Equal(t, loadedOrder.Number, originalOrder.Number)
	assert.Equal(t, loadedOrder.Accrual, originalOrder.Accrual)
	assert.Equal(t, loadedOrder.Status, originalOrder.Status)
	loadedWithdraw := loadedUser.Withdrawals[0]
	originalWithdraw := user.Withdrawals[0]
	assert.Equal(t, loadedWithdraw.Order, originalWithdraw.Order)
	assert.Equal(t, loadedWithdraw.Sum, originalWithdraw.Sum)
}
