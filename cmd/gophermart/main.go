package main

import (
	"net/http"
	"os"

	"github.com/Nickolasll/gomart/internal/application"
	"github.com/Nickolasll/gomart/internal/config"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/Nickolasll/gomart/internal/presentation"
	"github.com/sirupsen/logrus"

	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Получился слишком жирный main.go
// Но зато все инициализируется явно
// Возможно стоит инкапсулировать коннекшен gorm?

func main() {
	log := logrus.New()
	cfg := config.GetConfig()
	sqlDB, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		log.Info(err)
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Info(err)
		panic(err)
	}
	userAggregateRepository := infrastructure.UserAggregateRepository{DB: *db}
	orderRepository := infrastructure.OrderRepository{DB: *db}
	balanceRepository := infrastructure.BalanceRepository{DB: *db}
	withdrawRepository := infrastructure.WithdrawRepository{DB: *db}
	accrualClient := infrastructure.AccrualClient{URL: cfg.AccrualSystemURL}
	err = userAggregateRepository.Init()
	if err != nil {
		log.Info(err)
		panic(err)
	}
	jose := application.JOSEService{TokenExp: cfg.TokenExp, SecretKey: cfg.SecretKey}
	app := application.CreateApplication(
		userAggregateRepository,
		orderRepository,
		balanceRepository,
		withdrawRepository,
		accrualClient,
		jose,
		log,
	)
	mux := presentation.ChiFactory(&app, &jose, log)
	err = http.ListenAndServe(cfg.ServerEndpoint, mux)
	if err != nil {
		log.Info(err)
		panic(err)
	}
	os.Exit(app.ShutDown())
}
