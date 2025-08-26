# Service Overview API

## Overview
This project is a **Go microservice** that provides APIs to view and search organizational services.  
Each service has a name, description, and versions.  

The API allows you to:
- List services with search, sorting, and pagination
- Get details of a specific service
- Get available versions of a service  

All endpoints are secured with JWT authentication.

---

## Design and Trade-offs
- **Layers:** The code is split into handlers (HTTP layer), middleware (auth), and repository (data access). 
- **In-Memory vs SQLite:**  
  - Default: In-memory map (simple, fast, resets on restart).  
  - SQLite: Persistent storage. To enable, set `USE_SQLITE=true`.  
- **Auth:** API uses JWT tokens. First call `/login` with credentials to get a token, then use it in the `Authorization: Bearer <token>` header for all other endpoints.

---

## Endpoints
- `POST /login` – authenticate and get a JWT.  
- `GET /services` – list services. Supports:  
  - `?search=contact` – filter by name/description containing “contact”.  
  - `?sort=<id|name>&order=<asc|desc>` – sort results.  
  - `?limit=<n>&offset=<n>` – paginate results.  
- `GET /services/{id}` – get details of a service by ID.  
- `GET /services/{id}/versions` – get available versions of a service.  

All `/services` endpoints require a valid JWT.

---

## Running

1. Install [Go](https://golang.org/dl/).  

2. Set a JWT secret:  
   
   **For In-Memory DB:**
   ```bash
   export JWT_SECRET="your_32_character_secret_key_here"
   ```
   
   **For Persistent DB (SQLite):**
   ```bash
   export USE_SQLITE=true
   export JWT_SECRET="your_32_character_secret_key_here"
   ```

3. Start the server:
   ```bash
   go run main.go
   ```

4. Login (get JWT token):
   ```bash
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"password"}'
   ```

5. List services (use token from login):
   ```bash
   curl http://localhost:8080/services \
     -H "Authorization: Bearer <your_token_here>"
   ```

6. With search + sort by name ascending:
   ```bash
   curl -H "Authorization: Bearer <your_token>" \
     "http://localhost:8080/services?search=contact&sort=name&order=asc"
   ```

7. With search + pagination (limit 2, offset 0):
   ```bash
   curl -H "Authorization: Bearer <your_token>" \
     "http://localhost:8080/services?search=contact&limit=2&offset=0"
   ```
