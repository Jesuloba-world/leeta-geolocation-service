package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type HealthResponse struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	} `json:"body"`
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health Check",
		Description: "Check if the API is running and healthy",
		Tags:        []string{"Health"},
	}, h.HealthCheck)
}

func (h *HealthHandler) HealthCheck(ctx context.Context, input *struct{}) (*HealthResponse, error) {
	return &HealthResponse{
		Body: struct {
			Status string `json:"status" example:"ok"`
		}{
			Status: "ok",
		},
	}, nil
}