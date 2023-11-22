package domain

type AccrualOrderResponse struct {
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type AccrualClientInterface interface {
	GetOrderStatus(number string) (AccrualOrderResponse, error)
}
