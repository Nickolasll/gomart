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
	ch                     chan domain.Order
	log                    *logrus.Logger
	wg                     *sync.WaitGroup
}

func (w Worker) Serve() {
	order := <-w.ch
	defer w.wg.Done()
	processed, err := w.ProcessingOrderUseCase.Execute(order)
	if err != nil {
		if errors.Is(err, domain.ErrAccrualIsBusy) {
			w.wg.Add(1)
			w.ch <- order
			time.Sleep(1 * time.Second)
		} else {
			w.log.Info(err)
		}
	}
	if !processed {
		w.wg.Add(1)
		w.ch <- order
	}
}
