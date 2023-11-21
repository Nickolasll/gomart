package domain

import "errors"

type AccrualOrderResponse struct {
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type AccrualClientInterface interface {
	GetOrderStatus(number string) (AccrualOrderResponse, error)
}

var ErrDocumentNotFound = errors.New("requested document not found")
var ErrAccrualIsBusy = errors.New("accrual service is not ready")
