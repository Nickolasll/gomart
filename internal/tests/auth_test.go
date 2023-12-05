package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nickolasll/gomart/internal/application"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizationInvalidTokenValue(t *testing.T) {
	router, err := Init()
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
	router, err := Init()
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/api/user/orders", nil)
	req.Header.Add("Authorization", expiredToken)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}
