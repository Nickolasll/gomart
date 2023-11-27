package main

import (
	"net/http"
	"os"

	"github.com/Nickolasll/gomart/internal/application"
	"github.com/Nickolasll/gomart/internal/config"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/Nickolasll/gomart/internal/presentation"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	cfg := config.GetConfig()
	db, err := infrastructure.EstablishConnection(cfg.DatabaseURI)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	userAggregateRepository := infrastructure.UserAggregateRepository{DB: *db}
	orderRepository := infrastructure.OrderRepository{DB: *db}
	balanceRepository := infrastructure.BalanceRepository{DB: *db}
	withdrawRepository := infrastructure.WithdrawRepository{DB: *db}
	accrualClient := infrastructure.AccrualClient{URL: cfg.AccrualSystemURL}
	err = userAggregateRepository.Init()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		panic(err)
	}
	os.Exit(app.ShutDown())
}
