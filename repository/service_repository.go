package repository

import "service-api/models"

type ServiceRepository interface {
	List(searchTerm, sortBy, order string, limit, offset int) ([]*models.Service, error)
	GetByID(id int) (*models.Service, error)
	GetVersions(id int) ([]string, error)
}
