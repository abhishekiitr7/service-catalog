package repository

import (
	"errors"
	"service-api/models"
	"sort"
	"strings"
	"sync"
)

type MemoryServiceRepository struct {
	services map[int]*models.Service
	lock     sync.RWMutex
}

// NewMemoryServiceRepository initializes an in-memory store.
func NewMemoryServiceRepository() *MemoryServiceRepository {
	services := map[int]*models.Service{
		1: {ID: 1, Name: "Locate Us", Description: "Reach to us", Versions: []string{"1.0.0", "1.1.0"}},
		2: {ID: 2, Name: "Collect Money", Description: "receive money", Versions: []string{"1.0.0"}},
		3: {ID: 3, Name: "Contact Us", Description: "our contact channels for customer support", Versions: []string{"1.0.0", "1.2.0"}},
		4: {ID: 4, Name: "Notifications", Description: "updates to customers", Versions: []string{"2.0.0", "2.1.0"}},
		// Add more if needed
	}
	return &MemoryServiceRepository{services: services}
}

func (m *MemoryServiceRepository) List(searchTerm, sortBy, order string, limit, offset int) ([]*models.Service, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var list []*models.Service
	for _, svc := range m.services {
		if searchTerm != "" {
			// filter by name or description containing searchTerm
			if !strings.Contains(strings.ToLower(svc.Name), strings.ToLower(searchTerm)) &&
				!strings.Contains(strings.ToLower(svc.Description), strings.ToLower(searchTerm)) {
				continue
			}
		}
		list = append(list, svc)
	}

	// sort by name or id
	if sortBy == "name" {
		sort.Slice(list, func(i, j int) bool {
			if order == "desc" {
				return list[i].Name > list[j].Name
			}
			return list[i].Name < list[j].Name
		})
	} else {
		sort.Slice(list, func(i, j int) bool {
			if order == "desc" {
				return list[i].ID > list[j].ID
			}
			return list[i].ID < list[j].ID
		})
	}

	// apply limit/offset for pagination
	start := offset
	if start > len(list) {
		start = len(list)
	}
	end := start + limit
	if limit <= 0 || end > len(list) {
		end = len(list)
	}
	return list[start:end], nil
}

func (m *MemoryServiceRepository) GetByID(id int) (*models.Service, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if svc, ok := m.services[id]; ok {
		return svc, nil
	}
	return nil, errors.New("service not found")
}

func (m *MemoryServiceRepository) GetVersions(id int) ([]string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if svc, ok := m.services[id]; ok {
		vcopy := make([]string, len(svc.Versions))
		copy(vcopy, svc.Versions)
		return vcopy, nil
	}
	return nil, errors.New("service not found")
}
