package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	router, err := Init()
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
	router, err := Init()
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
	router, err := Init()
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
	router, err := Init()
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
	router, err := Init()
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
	router, err := Init()
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
