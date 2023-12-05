package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nickolasll/gomart/internal/infrastructure"
	"github.com/Nickolasll/gomart/internal/presentation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBalanceOk(t *testing.T) {
	router, err := Init()
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
	balanceResponse := presentation.BalanceResponse{}
	body, err := io.ReadAll(responseRecorder.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &balanceResponse)
	require.NoError(t, err)
	assert.Equal(t, balanceResponse.Current, float64(0))
	assert.Equal(t, balanceResponse.Withdrawn, float64(0))
}
