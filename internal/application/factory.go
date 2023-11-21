package application

import (
	"sync"

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
	ch       chan domain.Order
	wg       *sync.WaitGroup
}

func CreateApplication(DB gorm.DB, jose JOSEService, url string, log *logrus.Logger) Application {
	var wg sync.WaitGroup
	channel := make(chan domain.Order)
	userAggregateRepository := infrastructure.UserAggregateRepository{DB: DB}
	orderRepository := infrastructure.OrderRepository{DB: DB}
	balanceRepository := infrastructure.BalanceRepository{DB: DB}
	withdrawRepository := infrastructure.WithdrawRepository{DB: DB}
	accrualClient := infrastructure.AccrualClient{URL: url}
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
		wg:                      &wg,
	}
	getOrdersUseCase := GetOrders{
		orderRepository: orderRepository,
	}
	getBalanceUseCase := GetBalance{
		balanceRepository: balanceRepository,
		log:               log,
	}
	uploadWithdrawUseCase := UploadWithdraw{
		withdrawRepository:      withdrawRepository,
		userAggregateRepository: userAggregateRepository,
		log:                     log,
	}
	getWithdrawalsUseCase := GetWithdrawals{
		withdrawRepository: withdrawRepository,
	}
	useCases := UseCases{
		Registration:   registrationUseCase,
		Login:          loginUseCase,
		UploadOrder:    uploadOrderUseCase,
		GetOrders:      getOrdersUseCase,
		GetBalance:     getBalanceUseCase,
		UploadWithdraw: uploadWithdrawUseCase,
		GetWithdrawals: getWithdrawalsUseCase,
	}
	processingOrderUseCase := ProcessingOrder{
		userAggregateRepository: userAggregateRepository,
		accrualClient:           accrualClient,
	}
	worker := Worker{
		ProcessingOrderUseCase: processingOrderUseCase,
		ch:                     channel,
		log:                    log,
		wg:                     &wg,
	}
	go worker.Serve()
	return Application{
		UseCases: useCases,
		wg:       &wg,
		ch:       channel,
	}
}

func (a Application) ShutDown() {
	a.wg.Wait()
	close(a.ch)
}
