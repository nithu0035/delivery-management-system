# Delivery Management System (Go)

A minimal but complete Delivery Management System built for an assignment. Features:
- Go (GORM) + PostgreSQL for persistence
- Redis for caching & pub/sub (status updates)
- JWT authentication (Customer / Admin)
- Asynchronous order lifecycle simulation using goroutines
- Dockerfile + docker-compose for easy local run
- Automated tests (unit tests for core flows)

## Quick start (Docker)
Requires Docker & docker-compose.

```bash
git clone <repo>
cd delivery-management-system
docker-compose up --build
# Server will be on http://localhost:8080
```

## Run locally without Docker
- Start Postgres and Redis locally (or use docker-compose services).
- Edit `config/config.go` if needed for DB connection.
- `go run ./cmd/server`

## Project structure (important files)
- `cmd/server/main.go` - entrypoint
- `internal/{auth,users,orders,tracking}` - core packages
- `migrations/` - SQL for initial tables
- `docker-compose.yml` & `Dockerfile`

## Notes for evaluator
- Orders automatically progress (created -> dispatched -> in_transit -> delivered) every 5 seconds.
- Cancelled orders stop progressing.
- JWT tokens include role (`customer` or `admin`).

