package application

import (
	"errors"
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
	processed, err := w.ProcessingOrderUseCase.Execute(order)
	if errors.Is(err, domain.ErrAccrualIsBusy) {
		time.Sleep(1 * time.Second)
	} else {
		w.log.Info(err)
		processed = true
	}
	return processed
}

func (w Worker) Serve() {
	defer w.wg.Done()
	for order := range w.ch {
		for processed := false; !processed; processed = w.routine(order) {
		}
	}
}
