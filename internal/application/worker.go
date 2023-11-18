package application

import (
	"errors"
	"time"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	ProcessingOrderUseCase ProcessingOrder
	ch                     chan domain.Order
	log                    *logrus.Logger
}

func (w Worker) Serve() {
	order := <-w.ch
	processed, err := w.ProcessingOrderUseCase.Execute(order)
	if err != nil {
		if errors.Is(err, domain.ErrAccrualIsBusy) {
			time.Sleep(1 * time.Second)
			w.ch <- order
		} else {
			w.log.Info(err)
		}
	}
	if !processed {
		w.ch <- order
	}
}
