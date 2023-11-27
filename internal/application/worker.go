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

func (w Worker) Serve() {
	defer w.wg.Done()
	for order := range w.ch {
		processed := false
		for !processed {
			processed, err := w.ProcessingOrderUseCase.Execute(order)
			if err != nil {
				if errors.Is(err, domain.ErrAccrualIsBusy) {
					time.Sleep(1 * time.Second)
				} else {
					w.log.Info(err)
					processed = true
				}
			}
			if processed {
				break
			}
		}
	}
}
