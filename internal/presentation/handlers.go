package presentation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Nickolasll/gomart/internal/application"
	"github.com/google/uuid"
)

// Нужно ли везде проверять Content-Type?

type AuthenticatedHandler func(w http.ResponseWriter, r *http.Request, UserID uuid.UUID)

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload RegistrationPayload
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
	tokenString, err := app.UseCases.Registration.Execute(requestPayload.Login, requestPayload.Password)
	if err != nil {
		if errors.Is(err, application.ErrLoginAlreadyInUse) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload RegistrationPayload
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
	tokenString, err := app.UseCases.Login.Execute(requestPayload.Login, requestPayload.Password)
	if err != nil {
		if errors.Is(err, application.ErrLoginOrPasswordIsInvalid) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
}

func UploadOrderHandler(w http.ResponseWriter, r *http.Request, UserID uuid.UUID) {
	number, err := io.ReadAll(r.Body)
	if err != nil || len(number) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.UseCases.UploadOrder.Execute(UserID, string(number))
	if err != nil {
		if errors.Is(err, application.ErrNotValidNumber) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else if errors.Is(err, application.ErrOrderUploadedByThisUser) {
			w.WriteHeader(http.StatusOK)
		} else if errors.Is(err, application.ErrOrderUploadedByAnotherUser) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
