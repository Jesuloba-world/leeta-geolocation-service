package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/jesuloba-world/leeta-task/config"
	"github.com/jesuloba-world/leeta-task/internal/handlers"
	"github.com/jesuloba-world/leeta-task/internal/repository/memory"
	"github.com/jesuloba-world/leeta-task/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize repository
	repo := memory.NewInMemoryLocationRepository()

	// Initialize service
	locationService := service.NewLocationService(repo)

	// Initialize handlers
	locationHandler := handlers.NewLocationHandler(locationService)
	healthHandler := handlers.NewHealthHandler()

	// Create ServeMux
	mux := http.NewServeMux()

	// Create Huma API configuration
	config := huma.DefaultConfig("Location API", "1.0.0")
	config.Info.Description = "A RESTful API for managing geolocated stations with nearest location search capabilities"
	config.Info.Contact = &huma.Contact{
		Name: "Location API Team",
	}
	config.Servers = []*huma.Server{
		{URL: fmt.Sprintf("http://localhost:%s", cfg.Server.Port), Description: "Development server"},
	}

	// Create Huma API with humago adapter
	api := humago.New(mux, config)

	// Register all routes with Huma
	healthHandler.RegisterRoutes(api)
	locationHandler.RegisterRoutes(api)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info("Starting server", "port", cfg.Server.Port)
		slog.Info("API Documentation available", "url", fmt.Sprintf("http://localhost:%s/docs", cfg.Server.Port))
		slog.Info("OpenAPI JSON available", "url", fmt.Sprintf("http://localhost:%s/openapi.json", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited")
}
