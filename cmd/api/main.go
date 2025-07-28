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

	"github.com/jesuloba-world/leeta-task/internal/config"
	"github.com/jesuloba-world/leeta-task/internal/handlers"
	"github.com/jesuloba-world/leeta-task/internal/repository"
	"github.com/jesuloba-world/leeta-task/internal/service"
)

func main() {
	// Load configuration from environment
	cfg := config.LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize repository
	locationRepo, cleanup, err := repository.NewRepositoryFromConfig(cfg)
	if err != nil {
		slog.Error("Failed to initialize repository", "error", err)
		os.Exit(1)
	}

	slog.Info("Repository initialized", "type", cfg.Storage)

	// Initialize service
	locationService := service.NewLocationService(locationRepo)

	// Initialize handlers
	locationHandler := handlers.NewLocationHandler(locationService)
	healthHandler := handlers.NewHealthHandler()

	// Create ServeMux
	mux := http.NewServeMux()

	// Create Huma API configuration
	config := huma.DefaultConfig("Leeta Location API", "1.0.0")
	config.Info.Description = "A RESTful API for managing geolocated stations with nearest location search capabilities"
	config.Info.Contact = &huma.Contact{
		Name:  "Jesuloba John Abere",
		Email: "jesulobajohn@gmail.com",
	}
	config.Servers = []*huma.Server{
		{URL: fmt.Sprintf("http://localhost:%d", cfg.Server.Port), Description: "Development server"},
	}

	// Create Huma API with humago adapter
	api := humago.New(mux, config)

	// Register all routes with Huma
	healthHandler.RegisterRoutes(api)
	locationHandler.RegisterRoutes(api)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	go func() {
		slog.Info("Starting server", "port", cfg.Server.Port)
		slog.Info("API Documentation available", "url", fmt.Sprintf("http://localhost:%d/docs", cfg.Server.Port))
		slog.Info("OpenAPI JSON available", "url", fmt.Sprintf("http://localhost:%d/openapi.json", cfg.Server.Port))
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
	}

	// Cleanup database connection
	if err := cleanup(); err != nil {
		slog.Error("Failed to cleanup database connection", "error", err)
	}

	slog.Info("Server shutdown complete")
}
