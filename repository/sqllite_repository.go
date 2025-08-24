package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"service-api/models"
	"strings"
)

type SQLiteServiceRepository struct {
	db *sql.DB
}

func NewSQLiteServiceRepository(db *sql.DB) (*SQLiteServiceRepository, error) {
	sqLite := &SQLiteServiceRepository{db: db}
	err := sqLite.CreateIfRequired()
	if err != nil {
		return nil, err
	}
	return sqLite, nil
}

// CreateIfRequired creates the services table if not exists
func (r *SQLiteServiceRepository) CreateIfRequired() error {
	query := `CREATE TABLE IF NOT EXISTS services (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        versions TEXT
    );`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	// Insert data if table is empty
	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM services").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		seedServices := []models.Service{
			{ID: 1, Name: "Locate Us", Description: "Reach to us", Versions: []string{"1.0.0", "1.1.0"}},
			{ID: 2, Name: "Collect Money", Description: "receive money", Versions: []string{"1.0.0"}},
			{ID: 3, Name: "Contact Us", Description: "our contact channels for customer support", Versions: []string{"1.0.0", "1.2.0"}},
			{ID: 4, Name: "Notifications", Description: "updates to customers", Versions: []string{"2.0.0", "2.1.0"}},
		}
		for _, s := range seedServices {
			vers, _ := json.Marshal(s.Versions)
			_, err := r.db.Exec("INSERT INTO services (id, name, description, versions) VALUES (?, ?, ?, ?)",
				s.ID, s.Name, s.Description, string(vers))
			if err != nil {
				return fmt.Errorf("failed to seed: %v", err)
			}
		}
	}

	return nil
}

func (r *SQLiteServiceRepository) List(searchTerm, sortBy, order string, limit, offset int) ([]*models.Service, error) {
	if sortBy != "name" && sortBy != "id" {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	query := fmt.Sprintf("SELECT id, name, description, versions FROM services WHERE name LIKE ? OR description LIKE ? ORDER BY %s %s LIMIT ? OFFSET ?", sortBy, order)
	rows, err := r.db.Query(query, "%"+searchTerm+"%", "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*models.Service
	for rows.Next() {
		var s models.Service
		var versions string
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &versions); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(versions), &s.Versions)
		services = append(services, &s)
	}
	return services, nil
}

func (r *SQLiteServiceRepository) GetByID(id int) (*models.Service, error) {
	row := r.db.QueryRow("SELECT id, name, description, versions FROM services WHERE id = ?", id)
	var s models.Service
	var versions string
	if err := row.Scan(&s.ID, &s.Name, &s.Description, &versions); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("service not found")
		}
		return nil, err
	}
	json.Unmarshal([]byte(versions), &s.Versions)
	return &s, nil
}

func (r *SQLiteServiceRepository) GetVersions(id int) ([]string, error) {
	row := r.db.QueryRow("SELECT versions FROM services WHERE id = ?", id)
	var versions string
	if err := row.Scan(&versions); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("service not found")
		}
		return nil, err
	}
	var vers []string
	json.Unmarshal([]byte(versions), &vers)
	return vers, nil
}
