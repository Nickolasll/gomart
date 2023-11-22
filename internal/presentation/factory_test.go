package presentation

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Nickolasll/gomart/internal/application"
	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistrationBadRequest(t *testing.T) {
	tests := []struct {
		name string
		body []byte
	}{
		{
			name: "no password",
			body: []byte(`{"login": "no_password"}`),
		},
		{
			name: "no login",
			body: []byte(`{"password": "no_login"}`),
		},
		{
			name: "wrong fields",
			body: []byte(`{"field": "value"}`),
		},
		{
			name: "not a json",
			body: []byte(`not a json`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, err := ChiFactory()
			require.NoError(t, err)
			bodyReader := bytes.NewReader(tt.body)
			req := httptest.NewRequest("POST", "/api/user/register", bodyReader)
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestRegistrationSuccess(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	bodyReader := bytes.NewReader([]byte(`{"login": "` + uuid.NewString() + `", "password": "password"}`))
	req := httptest.NewRequest("POST", "/api/user/register", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Header().Get("Authorization"))
}

func TestRegistrationConflict(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	login := uuid.NewString()
	repo := infrastructure.UserAggregateRepository{DB: *db}
	repo.Create(login, "")
	bodyReader := bytes.NewReader([]byte(`{"login": "` + login + `", "password": "password"}`))
	req := httptest.NewRequest("POST", "/api/user/register", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusConflict, responseRecorder.Code)
}

func TestLoginBadRequest(t *testing.T) {
	tests := []struct {
		name string
		body []byte
	}{
		{
			name: "no password",
			body: []byte(`{"login": "no_password"}`),
		},
		{
			name: "no login",
			body: []byte(`{"password": "no_login"}`),
		},
		{
			name: "wrong fields",
			body: []byte(`{"field": "value"}`),
		},
		{
			name: "not a json",
			body: []byte(`not a json`),
		},
	}
	router, err := ChiFactory()
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := bytes.NewReader(tt.body)
			req := httptest.NewRequest("POST", "/api/user/login", bodyReader)
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestLoginNoUserUnauthorized(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	bodyReader := bytes.NewReader([]byte(`{"login": "` + uuid.NewString() + `", "password": "password"}`))
	req := httptest.NewRequest("POST", "/api/user/login", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestLoginSuccess(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	login := uuid.NewString()
	password := "password"
	repo := infrastructure.UserAggregateRepository{DB: *db}
	jose := application.JOSEService{TokenExp: 0, SecretKey: ""}
	repo.Create(login, jose.Hash(password))
	bodyReader := bytes.NewReader([]byte(`{"login": "` + login + `", "password": "` + password + `"}`))
	req := httptest.NewRequest("POST", "/api/user/login", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Header().Get("Authorization"))
}

func TestLoginWrongPasswordUnauthorized(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	login := uuid.NewString()
	repo := infrastructure.UserAggregateRepository{DB: *db}
	repo.Create(login, "password")
	bodyReader := bytes.NewReader([]byte(`{"login": "` + login + `", "password": "qwerty"}`))
	req := httptest.NewRequest("POST", "/api/user/login", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestAuthorizationInvalidTokenValue(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/user/orders", nil)
	req.Header.Add("Authorization", "invalid token value")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestAuthorizationExpiredToken(t *testing.T) {
	userID := uuid.New()
	expiredToken, err := application.JOSEService{TokenExp: 0, SecretKey: ""}.IssueToken(userID)
	require.NoError(t, err)
	router, err := ChiFactory()
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/user/orders", nil)
	req.Header.Add("Authorization", expiredToken)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestUploadOrderBadRequest(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	userID := uuid.New()
	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/user/orders", nil)
	req.Header.Add("Authorization", tokenValue)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func MockAuth() (*domain.UserAggregate, string, error) {
	var user domain.UserAggregate
	login := uuid.NewString()
	repo := infrastructure.UserAggregateRepository{DB: *db}
	user, err := repo.Create(login, "password")
	if err != nil {
		return nil, "", err
	}
	tokenValue, err := jose.IssueToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	return &user, tokenValue, nil
}

func TestUploadOrderUnprocessableEntity(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "not a number",
			body: "not a number",
		},
		{
			name: "invalid number",
			body: "4147203059780942",
		},
	}
	router, err := ChiFactory()
	require.NoError(t, err)
	userID := uuid.New()
	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.body)
			req := httptest.NewRequest("POST", "/api/user/orders", bodyReader)
			req.Header.Add("Authorization", tokenValue)
			req.Header.Add("Content-Type", "text/plain")
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusUnprocessableEntity, responseRecorder.Code)
		})
	}
}

func TestUploadOrderOtherUserUploaded(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	orderNumber := "9278923470"
	userID := uuid.New()
	user := domain.UserAggregate{ID: userID}
	repo := infrastructure.UserAggregateRepository{DB: *db}
	repo.Save(user)
	user, _ = user.AddOrder(orderNumber)
	repo.Save(user)

	tokenValue, err := jose.IssueToken(uuid.New())
	require.NoError(t, err)
	bodyReader := strings.NewReader(orderNumber)
	req := httptest.NewRequest("POST", "/api/user/orders", bodyReader)
	req.Header.Add("Authorization", tokenValue)
	req.Header.Add("Content-Type", "text/plain")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusConflict, responseRecorder.Code)
}

func TestUploadOrderThisUserUploaded(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	orderNumber := "12345678903"
	userID, err := uuid.Parse("66f0d9d1-8aca-4ce3-b5e8-e8a68f173716")
	require.NoError(t, err)
	user := domain.UserAggregate{ID: userID}
	repo := infrastructure.UserAggregateRepository{DB: *db}
	repo.Save(user)
	user, _ = user.AddOrder(orderNumber)
	repo.Save(user)

	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	bodyReader := strings.NewReader(orderNumber)
	req := httptest.NewRequest("POST", "/api/user/orders", bodyReader)
	req.Header.Add("Authorization", tokenValue)
	req.Header.Add("Content-Type", "text/plain")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestUploadOrderSuccess(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	orderNumber := "346436439"
	err = db.Delete(&domain.Order{Number: orderNumber}).Error
	require.NoError(t, err)
	userID := uuid.New()
	require.NoError(t, err)
	user := domain.UserAggregate{ID: userID}
	repo := infrastructure.UserAggregateRepository{DB: *db}
	repo.Save(user)

	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	bodyReader := strings.NewReader(orderNumber)
	req := httptest.NewRequest("POST", "/api/user/orders", bodyReader)
	req.Header.Add("Authorization", tokenValue)
	req.Header.Add("Content-Type", "text/plain")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusAccepted, responseRecorder.Code)
}

func TestGetOrdersNoContent(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	tokenValue, err := jose.IssueToken(uuid.New())
	require.NoError(t, err)
	req := httptest.NewRequest("GET", "/api/user/orders", nil)
	req.Header.Add("Authorization", tokenValue)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusNoContent, responseRecorder.Code)
}

func TestGetOrdersOk(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	err = db.Delete(&domain.Order{Number: "1"}).Error
	require.NoError(t, err)
	err = db.Delete(&domain.Order{Number: "2"}).Error
	require.NoError(t, err)
	userID := uuid.New()
	user := domain.UserAggregate{ID: userID}
	repo := infrastructure.UserAggregateRepository{DB: *db}
	repo.Save(user)
	order := domain.Order{
		UserAggregateID: userID,
		Number:          "1",
		UploadedAt:      time.Now(),
		Status:          domain.StatusProcessed,
		Accrual:         42050,
	}
	user.Orders = append(user.Orders, order)
	invalidOrder := domain.Order{
		UserAggregateID: userID,
		Number:          "2",
		UploadedAt:      time.Now(),
		Status:          domain.StatusInvalid,
		Accrual:         0,
	}
	user.Orders = append(user.Orders, invalidOrder)
	repo.Save(user)

	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	req := httptest.NewRequest("GET", "/api/user/orders", nil)
	req.Header.Add("Authorization", tokenValue)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestGetBalanceOk(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	repo := infrastructure.UserAggregateRepository{DB: *db}
	user, err := repo.Create(uuid.NewString(), "")
	require.NoError(t, err)
	tokenValue, err := jose.IssueToken(user.ID)
	require.NoError(t, err)
	req := httptest.NewRequest("GET", "/api/user/balance", nil)
	req.Header.Add("Authorization", tokenValue)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	balanceResponse := BalanceResponse{}
	body, err := io.ReadAll(responseRecorder.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &balanceResponse)
	require.NoError(t, err)
	assert.Equal(t, balanceResponse.Current, float64(0))
	assert.Equal(t, balanceResponse.Withdrawn, float64(0))
}

func TestUploadWithdrawBadRequest(t *testing.T) {
	tests := []struct {
		name string
		body []byte
	}{
		{
			name: "no number",
			body: []byte(`{"sum": "500.50"}`),
		},
		{
			name: "no sum",
			body: []byte(`{"order": "12345"}`),
		},
		{
			name: "wrong fields",
			body: []byte(`{"field": "value"}`),
		},
		{
			name: "not a json",
			body: []byte(`not a json`),
		},
	}
	router, err := ChiFactory()
	require.NoError(t, err)
	userID := uuid.New()
	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := bytes.NewReader(tt.body)
			req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bodyReader)
			req.Header.Add("Authorization", tokenValue)
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestUploadWithdrawOtherUserUploaded(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	number := "9278923470"
	repo := infrastructure.UserAggregateRepository{DB: *db}
	user, err := repo.Create(uuid.NewString(), "")
	require.NoError(t, err)
	user, err = user.AddWithdraw(number, 0)
	require.NoError(t, err)
	repo.Save(user)

	tokenValue, err := jose.IssueToken(uuid.New())
	require.NoError(t, err)
	bodyReader := bytes.NewReader([]byte(`{"order": "` + number + `", "sum": 500.50}`))
	req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bodyReader)
	req.Header.Add("Authorization", tokenValue)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusConflict, responseRecorder.Code)
}

func TestUploadWithdrawThisUserUploaded(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	number := "346436439"
	userID, err := uuid.Parse("b60f464c-0375-4c91-9deb-87cf0291d17c")
	require.NoError(t, err)
	repo := infrastructure.UserAggregateRepository{DB: *db}
	user := domain.UserAggregate{ID: userID}
	repo.Save(user)
	user.Balance = domain.Balance{Current: 0, Withdrawn: 0}
	user, err = user.AddWithdraw(number, 0)
	require.NoError(t, err)
	repo.Save(user)

	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	bodyReader := bytes.NewReader([]byte(`{"order": "` + number + `", "sum": 500.50}`))
	req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bodyReader)
	req.Header.Add("Authorization", tokenValue)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusConflict, responseRecorder.Code)
}

func TestUploadWithdrawUnprocessableEntity(t *testing.T) {
	tests := []struct {
		name string
		body []byte
	}{
		{
			name: "not a number",
			body: []byte(`{"order": "not a number", "sum": 500.50}`),
		},
		{
			name: "invalid number",
			body: []byte(`{"order": "4147203059780942", "sum": 500.50}`),
		},
	}
	router, err := ChiFactory()
	require.NoError(t, err)
	userID := uuid.New()
	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := bytes.NewReader(tt.body)
			req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bodyReader)
			req.Header.Add("Authorization", tokenValue)
			req.Header.Add("Content-Type", "application/json")
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusUnprocessableEntity, responseRecorder.Code)
		})
	}
}

func TestUploadWithdrawPaymentRequired(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	number := "12345678903"
	repo := infrastructure.UserAggregateRepository{DB: *db}
	user, err := repo.Create(uuid.NewString(), "")
	require.NoError(t, err)

	tokenValue, err := jose.IssueToken(user.ID)
	require.NoError(t, err)
	bodyReader := bytes.NewReader([]byte(`{"order": "` + number + `", "sum": 500.50}`))
	req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bodyReader)
	req.Header.Add("Authorization", tokenValue)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusPaymentRequired, responseRecorder.Code)
}

func TestUploadWithdrawSuccess(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	number := "346436439"
	err = db.Delete(&domain.Withdraw{Order: number}).Error
	require.NoError(t, err)
	userID := uuid.New()
	require.NoError(t, err)
	user := domain.UserAggregate{ID: userID, Login: userID.String()}
	repo := infrastructure.UserAggregateRepository{DB: *db}
	repo.Save(user)
	user.Balance.Current = 100000
	user.Balance.Withdrawn = 50000
	repo.Save(user)

	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	bodyReader := bytes.NewReader([]byte(`{"order": "` + number + `", "sum": 500}`))
	req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bodyReader)
	req.Header.Add("Authorization", tokenValue)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestGetWithdrawalsNoContent(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	tokenValue, err := jose.IssueToken(uuid.New())
	require.NoError(t, err)
	req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
	req.Header.Add("Authorization", tokenValue)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusNoContent, responseRecorder.Code)
}

func TestGetWithdrawalsOk(t *testing.T) {
	router, err := ChiFactory()
	require.NoError(t, err)
	err = db.Delete(&domain.Withdraw{Order: "1"}).Error
	require.NoError(t, err)
	userID := uuid.New()
	user := domain.UserAggregate{ID: userID}
	repo := infrastructure.UserAggregateRepository{DB: *db}
	repo.Save(user)
	withdraw := domain.Withdraw{
		UserAggregateID: userID,
		Order:           "1",
		ProcessedAt:     time.Now(),
		Sum:             42050,
	}
	user.Withdrawals = append(user.Withdrawals, withdraw)
	repo.Save(user)

	tokenValue, err := jose.IssueToken(userID)
	require.NoError(t, err)
	req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
	req.Header.Add("Authorization", tokenValue)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}
