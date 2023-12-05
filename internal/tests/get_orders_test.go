package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrdersNoContent(t *testing.T) {
	router, err := Init()
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
	router, err := Init()
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
