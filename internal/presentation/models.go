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
	Accrual    float32   `json:"accrual,omitempty"`
}

type BalanceResponse struct {
	Current  float32 `json:"current"`
	Withdraw float32 `json:"withdraw"`
}

type UploadWithdrawPayload struct {
	Number string  `json:"number"`
	Sum    float32 `json:"sum"`
}

type WithdrawalsResponse struct {
	Order       string    `json:"order"`
	ProcessedAt time.Time `json:"processed_at"`
	Sum         float32   `json:"sum,omitempty"`
}
