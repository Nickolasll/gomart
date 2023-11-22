package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

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
			router, err := Init()
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
	router, err := Init()
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
	router, err := Init()
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
