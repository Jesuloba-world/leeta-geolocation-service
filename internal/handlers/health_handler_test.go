package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupHealthTestAPI(t *testing.T) humatest.TestAPI {
	_, api := humatest.New(t, huma.DefaultConfig("Test API", "1.0.0"))

	healthHandler := NewHealthHandler()
	healthHandler.RegisterRoutes(api)

	return api
}

func TestHealthCheck(t *testing.T) {
	api := setupHealthTestAPI(t)

	resp := api.Get("/health")

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}