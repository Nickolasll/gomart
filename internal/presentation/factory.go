package presentation

import (
	"github.com/Nickolasll/gomart/internal/application"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var app *application.Application
var jose *application.JOSEService

func ChiFactory(
	App *application.Application,
	JOSE *application.JOSEService,
	Log *logrus.Logger,
) *chi.Mux {
	app = App
	log = Log
	jose = JOSE
	router := chi.NewRouter()
	router.Use(logging)
	router.Use(compress)

	authSubRouter := chi.NewRouter()
	authSubRouter.Post("/orders", auth(UploadOrderHandler))
	authSubRouter.Get("/orders", auth(GetOrdersHandler))
	authSubRouter.Get("/balance", auth(GetBalanceHandler))
	authSubRouter.Post("/balance/withdraw", auth(UploadWithdrawHandler))
	authSubRouter.Get("/withdrawals", auth(GetWithdrawalsHandler))

	router.Post("/api/user/register", RegistrationHandler)
	router.Post("/api/user/login", LoginHandler)
	router.Mount("/api/user", authSubRouter)
	return router
}
