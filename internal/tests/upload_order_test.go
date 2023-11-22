package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadOrderBadRequest(t *testing.T) {
	router, err := Init()
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
	router, err := Init()
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
	router, err := Init()
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
	router, err := Init()
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
	router, err := Init()
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
