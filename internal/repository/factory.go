package repository

import (
	"fmt"

	"github.com/jesuloba-world/leeta-task/internal/config"
	"github.com/jesuloba-world/leeta-task/internal/domain"
	"github.com/jesuloba-world/leeta-task/internal/repository/memory"
	"github.com/jesuloba-world/leeta-task/internal/repository/postgres"
)

const (
	MemoryRepository   = "memory"
	PostgresRepository = "postgres"
)

func NewRepositoryFromConfig(cfg config.Config) (domain.LocationRepository, func() error, error) {
	switch cfg.Storage {
	case MemoryRepository:
		return memory.NewInMemoryLocationRepository(), func() error { return nil }, nil
	case PostgresRepository:
		pgConfig := postgres.Config{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			DBName:   cfg.Database.DBName,
			SSLMode:  cfg.Database.SSLMode,
		}
		db, err := postgres.NewConnection(pgConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return postgres.NewPostgresLocationRepository(db), db.Close, nil
	default:
		return nil, nil, fmt.Errorf("unsupported repository type: %s", cfg.Storage)
	}
}
