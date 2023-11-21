package presentation

import "time"

type RegistrationPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OrderResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
	Accrual    float64   `json:"accrual,omitempty"`
}

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type UploadWithdrawPayload struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawalsResponse struct {
	Order       string    `json:"order"`
	ProcessedAt time.Time `json:"processed_at"`
	Sum         float64   `json:"sum,omitempty"`
}
