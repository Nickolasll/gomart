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

	// Я не придумал, как сделать rourer.Use(auth)
	// Потому что оно использует ServeHTTP(w ResponseWriter, r *Request)
	// И как пробросить UserID явно я не смог найти именно для chi
	authSubRouter := chi.NewRouter()
	authSubRouter.Post("/orders", auth(UploadOrderHandler))
	authSubRouter.Get("/orders", auth(GetOrdersHandler))
	authSubRouter.Get("/balance", auth(GetBalanceHandler))
	authSubRouter.Post("/balance/withdraw", auth(UploadWithdrawHandler))
	authSubRouter.Get("/withdrawals", auth(GetWithdrawalsHandler))

	router.Post("/api/user/register", RegistrationHandler)
	router.Post("/api/user/login", LoginHandler)
	router.Mount("/api/user", authSubRouter)
	return router, err
}

func AtExit() int {
	app.ShutDown()
	return 0
}
