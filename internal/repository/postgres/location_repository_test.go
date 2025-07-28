package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/jesuloba-world/leeta-task/internal/domain"
)

func setupTestContainer(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgis/postgis:17-3.5-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Enable PostGIS extension
	if _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS postgis"); err != nil {
		t.Fatalf("Failed to create PostGIS extension: %v", err)
	}

	// Create test table with PostGIS support
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS locations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			geom GEOGRAPHY(POINT, 4326),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	if _, err := db.Exec(createTableQuery); err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Create spatial index
	if _, err := db.Exec("CREATE INDEX IF NOT EXISTS idx_locations_geom ON locations USING GIST (geom)"); err != nil {
		t.Fatalf("Failed to create spatial index: %v", err)
	}

	// Create trigger to update geometry column
	triggerQuery := `
		CREATE OR REPLACE FUNCTION update_location_geom()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.geom = ST_Point(NEW.longitude, NEW.latitude)::geography;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		DROP TRIGGER IF EXISTS trigger_update_location_geom ON locations;
		CREATE TRIGGER trigger_update_location_geom
			BEFORE INSERT OR UPDATE ON locations
			FOR EACH ROW EXECUTE FUNCTION update_location_geom();
	`
	if _, err := db.Exec(triggerQuery); err != nil {
		t.Fatalf("Failed to create trigger: %v", err)
	}

	cleanup := func() {
		if _, err := db.Exec("DELETE FROM locations"); err != nil {
			t.Logf("Failed to clean up test data: %v", err)
		}
		db.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

func TestPostgresLocationRepository_Save(t *testing.T) {
	t.Run("successful save", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		location, err := domain.NewLocation("Test Location", 40.7128, -74.0060)
		if err != nil {
			t.Fatalf("Failed to create location: %v", err)
		}

		err = repo.Save(location)
		if err != nil {
			t.Fatalf("Failed to save location: %v", err)
		}

		if location.ID == "" {
			t.Error("Expected location ID to be set after save")
		}

		if location.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set after save")
		}
	})

	t.Run("duplicate name error", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		location1, _ := domain.NewLocation("Duplicate Location", 40.7128, -74.0060)
		location2, _ := domain.NewLocation("Duplicate Location", 41.8781, -87.6298)

		err := repo.Save(location1)
		if err != nil {
			t.Fatalf("Failed to save first location: %v", err)
		}

		err = repo.Save(location2)
		if err != domain.ErrLocationExists {
			t.Errorf("Expected ErrLocationExists, got: %v", err)
		}
	})
}

func TestPostgresLocationRepository_FindByName(t *testing.T) {
	t.Run("find existing location", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		location, _ := domain.NewLocation("Find Test Location", 40.7128, -74.0060)
		err := repo.Save(location)
		if err != nil {
			t.Fatalf("Failed to save location: %v", err)
		}

		found, err := repo.FindByName("Find Test Location")
		if err != nil {
			t.Fatalf("Failed to find location: %v", err)
		}

		if found.Name != "Find Test Location" {
			t.Errorf("Expected name 'Find Test Location', got: %s", found.Name)
		}
		if found.Latitude != 40.7128 {
			t.Errorf("Expected latitude 40.7128, got: %f", found.Latitude)
		}
		if found.Longitude != -74.0060 {
			t.Errorf("Expected longitude -74.0060, got: %f", found.Longitude)
		}
	})

	t.Run("location not found", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		_, err := repo.FindByName("Non-existent Location")
		if err != domain.ErrLocationNotFound {
			t.Errorf("Expected ErrLocationNotFound, got: %v", err)
		}
	})
}

func TestPostgresLocationRepository_FindByID(t *testing.T) {
	t.Run("find existing location by ID", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		location, _ := domain.NewLocation("ID Test Location", 40.7128, -74.0060)
		err := repo.Save(location)
		if err != nil {
			t.Fatalf("Failed to save location: %v", err)
		}

		found, err := repo.FindByID(location.ID)
		if err != nil {
			t.Fatalf("Failed to find location by ID: %v", err)
		}

		if found.ID != location.ID {
			t.Errorf("Expected ID %s, got: %s", location.ID, found.ID)
		}
		if found.Name != "ID Test Location" {
			t.Errorf("Expected name 'ID Test Location', got: %s", found.Name)
		}
	})

	t.Run("location not found by ID", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		_, err := repo.FindByID("999999")
		if err != domain.ErrLocationNotFound {
			t.Errorf("Expected ErrLocationNotFound, got: %v", err)
		}
	})
}

func TestPostgresLocationRepository_FindAll(t *testing.T) {
	t.Run("find all locations", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		locations := []*domain.Location{
			{Name: "Location 1", Latitude: 40.7128, Longitude: -74.0060, CreatedAt: time.Now()},
			{Name: "Location 2", Latitude: 41.8781, Longitude: -87.6298, CreatedAt: time.Now()},
			{Name: "Location 3", Latitude: 34.0522, Longitude: -118.2437, CreatedAt: time.Now()},
		}

		for _, loc := range locations {
			err := repo.Save(loc)
			if err != nil {
				t.Fatalf("Failed to save location %s: %v", loc.Name, err)
			}
		}

		found, err := repo.FindAll()
		if err != nil {
			t.Fatalf("Failed to find all locations: %v", err)
		}

		if len(found) != 3 {
			t.Errorf("Expected 3 locations, got: %d", len(found))
		}

		// Verify locations are ordered by ID
		for i := 1; i < len(found); i++ {
			if found[i-1].ID >= found[i].ID {
				t.Error("Locations should be ordered by ID")
			}
		}
	})

	t.Run("empty result", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		found, err := repo.FindAll()
		if err != nil {
			t.Fatalf("Failed to find all locations: %v", err)
		}

		if len(found) != 0 {
			t.Errorf("Expected 0 locations, got: %d", len(found))
		}
	})
}

func TestPostgresLocationRepository_Delete(t *testing.T) {
	t.Run("delete existing location", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		location, _ := domain.NewLocation("Delete Test Location", 40.7128, -74.0060)
		err := repo.Save(location)
		if err != nil {
			t.Fatalf("Failed to save location: %v", err)
		}

		err = repo.Delete(location.Name)
		if err != nil {
			t.Errorf("Failed to delete location: %v", err)
		}

		// Verify location is deleted
		_, err = repo.FindByName(location.Name)
		if err != domain.ErrLocationNotFound {
			t.Errorf("Expected ErrLocationNotFound after deletion, got: %v", err)
		}
	})

	t.Run("delete non-existent location", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		err := repo.Delete("Non-existent Location")
		if err != domain.ErrLocationNotFound {
			t.Errorf("Expected ErrLocationNotFound, got: %v", err)
		}
	})
}

func TestPostgresLocationRepository_FindNearest(t *testing.T) {
	t.Run("find nearest location", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		// Create test locations
		locations := []*domain.Location{
			{Name: "New York", Latitude: 40.7128, Longitude: -74.0060, CreatedAt: time.Now()},
			{Name: "Los Angeles", Latitude: 34.0522, Longitude: -118.2437, CreatedAt: time.Now()},
			{Name: "Chicago", Latitude: 41.8781, Longitude: -87.6298, CreatedAt: time.Now()},
			{Name: "Miami", Latitude: 25.7617, Longitude: -80.1918, CreatedAt: time.Now()},
		}

		for _, location := range locations {
			err := repo.Save(location)
			if err != nil {
				t.Fatalf("Failed to save location %s: %v", location.Name, err)
			}
		}

		// Test finding nearest to a point close to New York
		// Using coordinates slightly offset from New York
		nearestLocation, distance, err := repo.FindNearest(40.7500, -74.0000)
		if err != nil {
			t.Fatalf("Failed to find nearest location: %v", err)
		}

		if nearestLocation.Name != "New York" {
			t.Errorf("Expected nearest location to be 'New York', got '%s'", nearestLocation.Name)
		}

		if distance <= 0 {
			t.Errorf("Expected distance to be positive, got %f", distance)
		}
	})

	t.Run("no locations found", func(t *testing.T) {
		db, cleanup := setupTestContainer(t)
		defer cleanup()
		repo := NewPostgresLocationRepository(db)

		_, _, err := repo.FindNearest(40.7500, -74.0000)
		if err != domain.ErrLocationNotFound {
			t.Errorf("Expected ErrLocationNotFound when no locations exist, got: %v", err)
		}
	})
}

// Benchmark tests
func BenchmarkPostgresLocationRepository_Save(b *testing.B) {
	// Create a testing.T instance for setupTestContainer
	t := &testing.T{}
	db, cleanup := setupTestContainer(t)
	defer cleanup()

	repo := NewPostgresLocationRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		location, _ := domain.NewLocation(fmt.Sprintf("Benchmark Location %d", i), 40.7128, -74.0060)
		err := repo.Save(location)
		if err != nil {
			b.Fatalf("Failed to save location: %v", err)
		}
	}
}

func BenchmarkPostgresLocationRepository_FindByName(b *testing.B) {
	// Create a testing.T instance for setupTestContainer
	t := &testing.T{}
	db, cleanup := setupTestContainer(t)
	defer cleanup()

	repo := NewPostgresLocationRepository(db)

	// Setup test data
	location, _ := domain.NewLocation("Benchmark Location", 40.7128, -74.0060)
	err := repo.Save(location)
	if err != nil {
		b.Fatalf("Failed to save location: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.FindByName("Benchmark Location")
		if err != nil {
			b.Fatalf("Failed to find location: %v", err)
		}
	}
}
