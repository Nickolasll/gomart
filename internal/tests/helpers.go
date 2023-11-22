package tests

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Nickolasll/gomart/internal/application"
	"github.com/Nickolasll/gomart/internal/config"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/Nickolasll/gomart/internal/presentation"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var db *gorm.DB
var jose application.JOSEService

func Init() (*chi.Mux, error) {
	log := logrus.New()
	cfg := config.GetConfig()
	sqlDB, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	db, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	userAggregateRepository := infrastructure.UserAggregateRepository{DB: *db}
	orderRepository := infrastructure.OrderRepository{DB: *db}
	balanceRepository := infrastructure.BalanceRepository{DB: *db}
	withdrawRepository := infrastructure.WithdrawRepository{DB: *db}
	accrualClient := infrastructure.AccrualClient{URL: cfg.AccrualSystemURL}
	err = userAggregateRepository.Init()
	if err != nil {
		return nil, err
	}
	jose = application.JOSEService{TokenExp: cfg.TokenExp, SecretKey: cfg.SecretKey}
	app := application.CreateApplication(
		userAggregateRepository,
		orderRepository,
		balanceRepository,
		withdrawRepository,
		accrualClient,
		jose,
		log,
	)
	router := presentation.ChiFactory(&app, &jose, log)
	return router, err
}
