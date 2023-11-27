package application

import (
	"context"
	"time"

	"sync"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	ProcessingOrderUseCase ProcessingOrder
	ch                     <-chan domain.Order
	log                    *logrus.Logger
	wg                     *sync.WaitGroup
}

func (w Worker) routine(order domain.Order) bool {
	// processed, err := w.ProcessingOrderUseCase.Execute(order)
	// if errors.Is(err, domain.ErrAccrualIsBusy) {
	// 	time.Sleep(1 * time.Second)
	// } else {
	// 	w.log.Info(err)
	// 	processed = true
	// }
	// return processed
	time.Sleep(3 * time.Second)
	return false
}

func (w Worker) Serve() {
	defer w.wg.Done()
	for order := range w.ch {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		for processed := false; !processed; processed = w.routine(order) {
			select {
			case <-ctx.Done():
				w.log.Error("Processing order time out")
				cancel()
				continue
			default:
			}
		}
	}
}
