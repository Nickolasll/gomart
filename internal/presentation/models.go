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
	Accrual    string    `json:"accrual,omitempty"`
}

type BalanceResponse struct {
	Current  string `json:"current"`
	Withdraw string `json:"withdraw"`
}

type UploadWithdrawPayload struct {
	Number string `json:"number"`
	Sum    string `json:"sum"`
}

type WithdrawalsResponse struct {
	Order       string    `json:"order"`
	ProcessedAt time.Time `json:"processed_at"`
	Sum         string    `json:"sum,omitempty"`
}
