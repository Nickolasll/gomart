package presentation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Nickolasll/gomart/internal/application"
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/google/uuid"
)

type AuthenticatedHandler func(w http.ResponseWriter, r *http.Request, UserID uuid.UUID)

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload RegistrationPayload
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &requestPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if requestPayload.Login == "" || requestPayload.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenString, err := app.Registration.Execute(requestPayload.Login, requestPayload.Password)
	if err != nil {
		if errors.Is(err, application.ErrLoginAlreadyInUse) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}
		return
	}

	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload RegistrationPayload
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &requestPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if requestPayload.Login == "" || requestPayload.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenString, err := app.Login.Execute(requestPayload.Login, requestPayload.Password)
	if err != nil {
		if errors.Is(err, application.ErrLoginOrPasswordIsInvalid) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}
		return
	}

	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
}

func UploadOrderHandler(w http.ResponseWriter, r *http.Request, UserID uuid.UUID) {
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	number, err := io.ReadAll(r.Body)
	if err != nil || len(number) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.UploadOrder.Execute(UserID, string(number))
	if err != nil {
		switch {
		case errors.Is(err, application.ErrNotValidNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, application.ErrUploadedByThisUser):
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, application.ErrUploadedByAnotherUser):
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func GetOrdersHandler(w http.ResponseWriter, r *http.Request, UserID uuid.UUID) {
	w.Header().Set("Content-Type", "application/json")
	orders, err := app.GetOrders.Execute(UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ordersResponse := []OrderResponse{}
	for _, order := range orders {

		orderResponse := OrderResponse{
			Number:     order.Number,
			Status:     order.Status,
			UploadedAt: order.UploadedAt,
		}
		if orderResponse.Status == domain.StatusProcessed {
			orderResponse.Accrual = order.AccrualToFloat()
		}
		ordersResponse = append(ordersResponse, orderResponse)
	}
	resp, err := json.Marshal(ordersResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request, UserID uuid.UUID) {
	w.Header().Set("Content-Type", "application/json")
	balance, err := app.GetBalance.Execute(UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	balanceResponse := BalanceResponse{
		Current:   balance.CurrentToFloat(),
		Withdrawn: balance.WithdrawnToFloat(),
	}
	resp, err := json.Marshal(balanceResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func UploadWithdrawHandler(w http.ResponseWriter, r *http.Request, UserID uuid.UUID) {
	var requestPayload UploadWithdrawPayload
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &requestPayload)
	if err != nil || requestPayload.Order == "" || requestPayload.Sum == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.UploadWithdraw.Execute(UserID, requestPayload.Order, requestPayload.Sum)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrNotValidNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, application.ErrUploadedByAnotherUser):
			w.WriteHeader(http.StatusConflict)
		case errors.Is(err, domain.ErrInsufficientFunds):
			w.WriteHeader(http.StatusPaymentRequired)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request, UserID uuid.UUID) {
	w.Header().Set("Content-Type", "application/json")
	withdrawals, err := app.GetWithdrawals.Execute(UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	withdrawalsResponse := []WithdrawalsResponse{}
	for _, withdraw := range withdrawals {

		withdrawResponse := WithdrawalsResponse{
			Order:       withdraw.Order,
			Sum:         withdraw.SumToFloat(),
			ProcessedAt: withdraw.ProcessedAt,
		}

		withdrawalsResponse = append(withdrawalsResponse, withdrawResponse)
	}
	resp, err := json.Marshal(withdrawalsResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
