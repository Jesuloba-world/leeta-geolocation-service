package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/jesuloba-world/leeta-task/internal/domain"
)

type PostgresLocationRepository struct {
	db *sql.DB
}

func NewPostgresLocationRepository(db *sql.DB) *PostgresLocationRepository {
	return &PostgresLocationRepository{db: db}
}

func (r *PostgresLocationRepository) Save(location *domain.Location) error {
	existingLocation, err := r.FindByName(location.Name)
	if err == nil && existingLocation != nil {
		return domain.ErrLocationExists
	}

	query := `INSERT INTO locations (name, latitude, longitude) 
			 VALUES ($1, $2, $3) 
			 RETURNING id, created_at`

	var id int
	err = r.db.QueryRow(query, location.Name, location.Latitude, location.Longitude).Scan(&id, &location.CreatedAt)
	if err != nil {
		return err
	}

	location.ID = fmt.Sprintf("%d", id)
	return nil
}

func (r *PostgresLocationRepository) FindByName(name string) (*domain.Location, error) {
	query := `SELECT id, name, latitude, longitude, created_at 
			 FROM locations 
			 WHERE name = $1`

	var location domain.Location
	var id int
	err := r.db.QueryRow(query, name).Scan(
		&id,
		&location.Name,
		&location.Latitude,
		&location.Longitude,
		&location.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrLocationNotFound
		}
		return nil, err
	}

	location.ID = fmt.Sprintf("%d", id)
	return &location, nil
}

func (r *PostgresLocationRepository) FindByID(id string) (*domain.Location, error) {
	query := `SELECT id, name, latitude, longitude, created_at 
			 FROM locations 
			 WHERE id = $1`

	var location domain.Location
	var dbID int
	err := r.db.QueryRow(query, id).Scan(
		&dbID,
		&location.Name,
		&location.Latitude,
		&location.Longitude,
		&location.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrLocationNotFound
		}
		return nil, err
	}

	location.ID = fmt.Sprintf("%d", dbID)
	return &location, nil
}

func (r *PostgresLocationRepository) FindAll() ([]*domain.Location, error) {
	query := `SELECT id, name, latitude, longitude, created_at 
			 FROM locations 
			 ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := []*domain.Location{}
	for rows.Next() {
		var location domain.Location
		var id int
		err = rows.Scan(
			&id,
			&location.Name,
			&location.Latitude,
			&location.Longitude,
			&location.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		location.ID = fmt.Sprintf("%d", id)
		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *PostgresLocationRepository) Delete(name string) error {
	query := `DELETE FROM locations WHERE name = $1`

	result, err := r.db.Exec(query, name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrLocationNotFound
	}

	return nil
}

func (r *PostgresLocationRepository) FindNearest(latitude, longitude float64) (*domain.Location, float64, error) {
	query := `SELECT id, name, latitude, longitude, created_at,
				 ST_Distance(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) as distance
			  FROM locations 
			  ORDER BY geom <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography 
			  LIMIT 1`

	var location domain.Location
	var id int
	var distance float64
	err := r.db.QueryRow(query, longitude, latitude).Scan(
		&id,
		&location.Name,
		&location.Latitude,
		&location.Longitude,
		&location.CreatedAt,
		&distance,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, domain.ErrLocationNotFound
		}
		return nil, 0, err
	}

	location.ID = fmt.Sprintf("%d", id)
	return &location, distance, nil
}
