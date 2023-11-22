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

func TestGetWithdrawalsNoContent(t *testing.T) {
	router, err := Init()
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
	router, err := Init()
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
