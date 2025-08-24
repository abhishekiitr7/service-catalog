package handlers

import (
	"encoding/json"
	"net/http"
	"service-api/repository"
	"strconv"

	"github.com/gorilla/mux"
)

// ListServicesHandler returns a handler that lists services.
func ListServicesHandler(repo repository.ServiceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		q := r.URL.Query()
		search := q.Get("search")
		sortBy := q.Get("sort")
		if sortBy == "" {
			sortBy = "id"
		}
		order := q.Get("order")
		if order != "desc" {
			order = "asc"
		}

		limit := 10
		if l := q.Get("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil {
				limit = v
			}
		}
		offset := 0
		if o := q.Get("offset"); o != "" {
			if v, err := strconv.Atoi(o); err == nil {
				offset = v
			}
		}

		services, err := repo.List(search, sortBy, order, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
	}
}

// GetServiceHandler returns a handler that gets a service by ID.
func GetServiceHandler(repo repository.ServiceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch ID from URL - Path parameter
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		svc, err := repo.GetByID(id)
		if err != nil {
			http.Error(w, "service not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(svc)
	}
}

// GetServiceVersionsHandler returns a handler that lists versions of a service.
func GetServiceVersionsHandler(repo repository.ServiceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch ID from URL - Path parameter
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		versions, err := repo.GetVersions(id)
		if err != nil {
			http.Error(w, "service not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{"versions": versions})
	}
}
