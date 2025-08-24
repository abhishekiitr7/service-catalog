package repository

import (
	"database/sql"
	"reflect"
	"testing"

	"service-api/models"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	return db
}

func TestSQLiteServiceRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo, err := NewSQLiteServiceRepository(db)
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}

	tests := []struct {
		name       string
		gotFunc    func() (interface{}, error)
		want       interface{}
		shouldFail bool
	}{
		{
			name: "List services",
			gotFunc: func() (interface{}, error) {
				return repo.List("", "id", "asc", 10, 0)
			},
			want: []*models.Service{
				{ID: 1, Name: "Locate Us", Description: "Reach to us", Versions: []string{"1.0.0", "1.1.0"}},
				{ID: 2, Name: "Collect Money", Description: "receive money", Versions: []string{"1.0.0"}},
				{ID: 3, Name: "Contact Us", Description: "our contact channels for customer support", Versions: []string{"1.0.0", "1.2.0"}},
				{ID: 4, Name: "Notifications", Description: "updates to customers", Versions: []string{"2.0.0", "2.1.0"}},
			},
			shouldFail: false,
		},
		{
			name: "Get by valid ID",
			gotFunc: func() (interface{}, error) {
				return repo.GetByID(1)
			},
			want:       &models.Service{ID: 1, Name: "Locate Us", Description: "Reach to us", Versions: []string{"1.0.0", "1.1.0"}},
			shouldFail: false,
		},
		{
			name: "Get by invalid ID",
			gotFunc: func() (interface{}, error) {
				return repo.GetByID(999)
			},
			want:       nil,
			shouldFail: true,
		},
		{
			name: "Get versions",
			gotFunc: func() (interface{}, error) {
				return repo.GetVersions(1)
			},
			want:       []string{"1.0.0", "1.1.0"},
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.gotFunc()

			if tt.shouldFail {
				if err == nil {
					t.Errorf("expected error, got success with %+v", got)
				}
				return
			}

			if err != nil {
				t.Errorf("expected success, got error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
