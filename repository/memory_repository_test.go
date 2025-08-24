package repository

import (
	"reflect"
	"service-api/models"
	"sync"
	"testing"
)

func TestMemoryServiceRepository_List(t *testing.T) {
	type fields struct {
		services map[int]*models.Service
		// lock     sync.RWMutex
	}

	testFields := fields{
		services: map[int]*models.Service{
			1: {ID: 1, Name: "Locate Us", Description: "Reach to Us", Versions: []string{"1.0.0", "1.1.0"}},
			2: {ID: 2, Name: "Collect Money", Description: "recieve money", Versions: []string{"1.0.0"}},
			3: {ID: 3, Name: "Contact Us", Description: "our contact channels for customer support", Versions: []string{"1.0.0", "1.2.0"}},
			4: {ID: 4, Name: "Notifications", Description: "updates to customers", Versions: []string{"2.0.0", "2.1.0"}},
		},
		// lock: sync.RWMutex{},
	}

	tests := []struct {
		name       string
		fields     fields
		want       []*models.Service
		searchTerm string
		sortBy     string
		order      string
		limit      int
		offset     int
		wantErr    bool
	}{
		{
			name:       "List all services sorted by name ascending",
			fields:     testFields,
			searchTerm: "contact",
			sortBy:     "name",
			wantErr:    false,
			order:      "desc",
			want: []*models.Service{
				{ID: 3, Name: "Contact Us", Description: "our contact channels for customer support", Versions: []string{"1.0.0", "1.2.0"}},
			},
		},
		{
			name:       "List all services sorted by ID ascending",
			fields:     testFields,
			searchTerm: "contact",
			sortBy:     "id",
			wantErr:    false,
			order:      "desc",
			want: []*models.Service{
				{ID: 3, Name: "Contact Us", Description: "our contact channels for customer support", Versions: []string{"1.0.0", "1.2.0"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryServiceRepository{
				services: tt.fields.services,
				lock:     sync.RWMutex{},
			}
			got, err := m.List(tt.searchTerm, tt.sortBy, tt.order, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryServiceRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryServiceRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryServiceRepository_GetByID(t *testing.T) {
	type fields struct {
		services map[int]*models.Service
	}

	testFields := fields{
		services: map[int]*models.Service{
			1: {ID: 1, Name: "Locate Us", Description: "Reach to Us", Versions: []string{"1.0.0", "1.1.0"}},
			2: {ID: 2, Name: "Collect Money", Description: "recieve money", Versions: []string{"1.0.0"}},
			3: {ID: 3, Name: "Contact Us", Description: "our contact channels for customer support", Versions: []string{"1.0.0", "1.2.0"}},
			4: {ID: 4, Name: "Notifications", Description: "updates to customers", Versions: []string{"2.0.0", "2.1.0"}},
		},
	}

	tests := []struct {
		name    string
		fields  fields
		want    *models.Service
		id      int
		wantErr bool
	}{
		{
			name:    "Get service by valid ID",
			fields:  testFields,
			id:      2,
			wantErr: false,
			want: &models.Service{
				ID:          2,
				Name:        "Collect Money",
				Description: "recieve money",
				Versions:    []string{"1.0.0"},
			},
		},
		{
			name:    "Get service by invalid ID",
			fields:  testFields,
			id:      -1,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryServiceRepository{
				services: tt.fields.services,
				lock:     sync.RWMutex{},
			}
			got, err := m.GetByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryServiceRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryServiceRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
