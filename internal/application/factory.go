package application

import (
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Я хотел сделать пакет use_cases, но чтобы проект выглядел чище, но что-то пошло не так из-за exceptions
// Это точно не исключения предметной области, а именно исключения приложения, но я хотел бы более
// удобную структуру

type Application struct {
	UseCases UseCases
}

func CreateApplication(DB gorm.DB, jose JOSEService, url string, log *logrus.Logger) Application {
	channel := make(chan domain.Order)
	userAggregateRepository := infrastructure.UserAggregateRepository{DB: DB}
	orderRepository := infrastructure.OrderRepository{DB: DB}
	userAggregateRepository.Init()
	registrationUseCase := Registration{
		userAggregateRepository: userAggregateRepository,
		jose:                    jose,
	}
	loginUseCase := Login{
		userAggregateRepository: userAggregateRepository,
		jose:                    jose,
	}
	uploadOrderUseCase := UploadOrder{
		userAggregateRepository: userAggregateRepository,
		orderRepository:         orderRepository,
		ch:                      channel,
	}
	useCases := UseCases{
		Registration: registrationUseCase,
		Login:        loginUseCase,
		UploadOrder:  uploadOrderUseCase,
	}
	processingOrderUseCase := ProcessingOrder{
		userAggregateRepository: userAggregateRepository,
		accrualClient:           infrastructure.AccrualClient{URL: url},
	}
	worker := Worker{
		ProcessingOrderUseCase: processingOrderUseCase,
		ch:                     channel,
		log:                    log,
	}
	go worker.Serve()
	return Application{UseCases: useCases}
}
