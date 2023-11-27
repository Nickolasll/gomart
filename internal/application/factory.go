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
		log:                     log,
	}
	channel := make(chan domain.Order)
	worker := Worker{
		ProcessingOrderUseCase: processingOrderUseCase,
		ch:                     channel,
		log:                    log,
		wg:                     &wg,
	}
	wg.Add(1)
	go worker.Serve()
	registrationUseCase := registration{
		userAggregateRepository: userAggregateRepository,
		jose:                    jose,
		log:                     log,
	}
	loginUseCase := login{
		userAggregateRepository: userAggregateRepository,
		jose:                    jose,
		log:                     log,
	}
	uploadOrderUseCase := uploadOrder{
		userAggregateRepository: userAggregateRepository,
		orderRepository:         orderRepository,
		ch:                      channel,
		log:                     log,
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
		log:                     log,
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
	close(a.channel)
	a.waitGroup.Wait()
	return 0
}
