package presentation

import (
	"github.com/Nickolasll/gomart/internal/application"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	jose = application.JOSEService{TokenExp: TokenExp, SecretKey: SecretKey}
	app = application.CreateApplication(*db, jose, *AccrualSystemURL, log)
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
	return router, err
}
