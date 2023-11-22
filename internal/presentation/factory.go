package presentation

import (
	"github.com/Nickolasll/gomart/internal/application"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var log = logrus.New()
var app application.Application
var db *gorm.DB
var jose application.JOSEService

func ChiFactory() (*chi.Mux, error) {
	sqlDB, err := sql.Open("pgx", *DatabaseURI)
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
	userAggregateRepository.Init()
	accrualClient := infrastructure.AccrualClient{URL: *AccrualSystemURL}
	jose = application.JOSEService{TokenExp: TokenExp, SecretKey: SecretKey}
	app = application.CreateApplication(
		userAggregateRepository,
		orderRepository,
		balanceRepository,
		withdrawRepository,
		accrualClient,
		jose,
		log,
	)
	router := chi.NewRouter()
	router.Use(logging)
	router.Use(compress)
	// Думаю, как тут лучше разбить на саброутеры
	// И вообще стоит ли тут использовать саброутеры?
	// Возможно саброутер, использующий auth middleware
	// Посмотрим ближе к готовности
	router.Post("/api/user/register", RegistrationHandler)
	router.Post("/api/user/login", LoginHandler)
	router.Post("/api/user/orders", auth(UploadOrderHandler))
	router.Get("/api/user/orders", auth(GetOrdersHandler))
	router.Get("/api/user/balance", auth(GetBalanceHandler))
	router.Post("/api/user/balance/withdraw", auth(UploadWithdrawHandler))
	router.Get("/api/user/withdrawals", auth(GetWithdrawalsHandler))
	return router, err
}

func AtExit() int {
	app.ShutDown()
	return 0
}
