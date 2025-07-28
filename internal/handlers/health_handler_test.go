package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
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

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}