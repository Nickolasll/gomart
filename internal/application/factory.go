package application

import (
	"sync"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Registration   registration
	Login          login
	UploadOrder    uploadOrder
	GetOrders      getOrders
	GetBalance     getBalance
	UploadWithdraw uploadWithdraw
	GetWithdrawals getWithdrawals
	channel        chan domain.Order
	waitGroup      *sync.WaitGroup
}

func CreateApplication(
	userAggregateRepository domain.UserAggregateRepositoryInterface,
	orderRepository domain.OrderRepositoryInterface,
	balanceRepository domain.BalanceRepositoryInterface,
	withdrawRepository domain.WithdrawRepositoryInterface,
	accrualClient domain.AccrualClientInterface,
	jose JOSEService,
	log *logrus.Logger,
) Application {
	var wg sync.WaitGroup
	processingOrderUseCase := ProcessingOrder{
		userAggregateRepository: userAggregateRepository,
		accrualClient:           accrualClient,
	}
	channel := make(chan domain.Order)
	worker := Worker{
		ProcessingOrderUseCase: processingOrderUseCase,
		ch:                     channel,
		log:                    log,
		wg:                     &wg,
	}
	go worker.Serve()
	registrationUseCase := registration{
		userAggregateRepository: userAggregateRepository,
		jose:                    jose,
	}
	loginUseCase := login{
		userAggregateRepository: userAggregateRepository,
		jose:                    jose,
	}
	uploadOrderUseCase := uploadOrder{
		userAggregateRepository: userAggregateRepository,
		orderRepository:         orderRepository,
		ch:                      channel,
		wg:                      &wg,
	}
	getOrdersUseCase := getOrders{
		orderRepository: orderRepository,
	}
	getBalanceUseCase := getBalance{
		balanceRepository: balanceRepository,
	}
	uploadWithdrawUseCase := uploadWithdraw{
		withdrawRepository:      withdrawRepository,
		userAggregateRepository: userAggregateRepository,
	}
	getWithdrawalsUseCase := getWithdrawals{
		withdrawRepository: withdrawRepository,
	}
	return Application{
		Registration:   registrationUseCase,
		Login:          loginUseCase,
		UploadOrder:    uploadOrderUseCase,
		GetOrders:      getOrdersUseCase,
		GetBalance:     getBalanceUseCase,
		UploadWithdraw: uploadWithdrawUseCase,
		GetWithdrawals: getWithdrawalsUseCase,
		waitGroup:      &wg,
		channel:        channel,
	}
}

func (a Application) ShutDown() int {
	a.waitGroup.Wait()
	close(a.channel)
	return 0
}
