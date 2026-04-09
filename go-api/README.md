# Go Task Management API

A Go reimplementation of the [PHP Symfony Task Management REST API](../), built for 1:1 response parity. Both APIs run in Docker with MySQL 5.7 and produce identical JSON responses for the same inputs.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/task` | Create a new task |
| `GET` | `/api/task` | List all tasks (filter: `areDone`, `when`) |
| `GET` | `/api/task/{taskId}` | Get a specific task |
| `PUT` | `/api/task` | Update an existing task |
| `PUT` | `/api/task/{taskId}/complete` | Mark a task as completed |
| `DELETE` | `/api/task/{taskId}` | Delete a task |

## Tech Stack

- **Go 1.22** with [chi](https://github.com/go-chi/chi) router
- **MySQL 5.7** via [sqlx](https://github.com/jmoiron/sqlx) + [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- **Docker** for all services (API, database, tests)
- **testify** for test assertions

## Quick Start

```bash
# Build and run the Go API + MySQL
docker compose up -d --build

# The API is available at http://localhost:8080
curl http://localhost:8080/api/task
```

## Project Structure

```
go-api/
├── cmd/server/main.go              # Entry point, DI wiring
├── internal/
│   ├── config/                      # Environment variable loading
│   ├── database/                    # MySQL connection + migration runner
│   ├── handler/                     # HTTP handlers (1:1 with PHP controllers)
│   │   ├── task_handler.go          # Route handlers
│   │   └── response.go             # Request/response DTOs, PHP DateTime format
│   ├── model/                       # Domain model (Task)
│   ├── repository/                  # Database access layer
│   └── service/                     # Business logic + validation
├── migrations/                      # SQL migrations (identical to PHP)
├── docker/
│   ├── Dockerfile                   # Multi-stage production build
│   └── Dockerfile.test              # Test runner image
├── docker-compose.yml               # Standalone: Go API + MySQL
├── docker-compose.e2e.yml           # E2E tests: test runner + MySQL
├── docker-compose.test.yml          # Parity tests: PHP + Go + MySQL
├── tests/
│   ├── e2e/                         # Go endpoint tests (12 tests)
│   ├── parity/                      # Cross-parity tests (PHP vs Go)
│   └── fixtures/                    # Test data (route inventory, parity matrix)
├── Makefile
└── REPORT.md                        # Generated parity test report
```

## Running Tests

### All Tests

```bash
make test-all
```

This runs all three test suites in sequence: Go E2E, PHP, and Parity.

### Go E2E Tests

Runs 12 endpoint tests against the Go API with a real MySQL database in Docker:

```bash
make test-e2e
```

These tests mirror the PHP functional tests 1:1:
- `TestGetTask` / `TestGetTaskNotFound`
- `TestAddTask`
- `TestUpdateTask` / `TestUpdateTaskNotFound` / `TestUpdateTaskInvalidData` / `TestUpdateInvalidJson`
- `TestGetTasks`
- `TestCompleteTask` / `TestCompleteTaskNotFound`
- `TestDeleteTask` / `TestDeleteTaskNotFound`

### PHP Tests

Runs the existing PHP unit (24) and functional (12) test suites inside Docker:

```bash
make test-php
```

---

## Parity Testing

Parity testing is the core proof that the Go API is a correct migration of the PHP API. It sends **identical HTTP requests** to both APIs (which share the same MySQL database) and asserts that their responses match.

### What Parity Tests Verify

For each of the 12 test cases:
1. **Status codes** are identical (e.g., both return 200, or both return 404)
2. **Content-Type headers** both indicate `application/json`
3. **Response bodies** are structurally identical (after normalizing timestamps and auto-increment IDs)

A **coverage completeness meta-test** verifies that every route in the API has at least one parity test case.

### How to Run Parity Tests

#### Option 1: Using Make (recommended)

```bash
make test-parity
```

This will:
1. Build and start MySQL, PHP API (nginx + php-fpm), and Go API containers
2. Wait for all services to be healthy
3. Run the parity test suite
4. Generate `REPORT.md` with detailed request/response data
5. Tear down all containers

#### Option 2: Using Docker Compose directly

```bash
# Start all services
docker compose -f docker-compose.test.yml up -d --build

# Wait for services to be ready (MySQL health check + app startup)
sleep 20

# Run the parity tests
docker compose -f docker-compose.test.yml run --rm test-runner

# View the generated report
cat REPORT.md

# Clean up
docker compose -f docker-compose.test.yml down -v
```

#### Option 3: Manual setup (host networking)

If Docker bridge networking is unavailable (e.g., in some cloud VMs), you can run services individually with `--network host`:

```bash
# 1. Start MySQL on a custom port
docker run -d --name parity-mysql --network host \
  -e MYSQL_DATABASE=task-list \
  -e MYSQL_USER=secretuser \
  -e MYSQL_PASSWORD=thisisasupersecretpassworddontyouthink \
  -e MYSQL_ROOT_PASSWORD=iamgroot \
  -e MYSQL_TCP_PORT=3308 \
  mysql:5.7

# 2. Wait for MySQL
until docker exec parity-mysql mysqladmin ping -h 127.0.0.1 --port=3308 --silent 2>/dev/null; do sleep 2; done

# 3. Start Go API
docker run -d --name parity-go-api --network host \
  -e DB_SERVER=127.0.0.1 -e DB_PORT=3308 -e DB_NAME=task-list \
  -e DB_USER=secretuser -e DB_PASSWORD=thisisasupersecretpassworddontyouthink \
  -e APP_PORT=8082 \
  go-api:latest

# 4. Start PHP (php-fpm + nginx with host-mode config)
docker run -d --name parity-php-fpm --network host \
  -e DATABASE_URL="mysql://secretuser:thisisasupersecretpassworddontyouthink@127.0.0.1:3308/task-list?serverVersion=5.7" \
  tasklistfpm:latest

docker run -d --name parity-nginx --network host \
  -v $(pwd)/../docker/nginx/conf.d/parity.conf:/etc/nginx/conf.d/default.conf:ro \
  -v $(pwd)/..:/var/www/task-list:ro \
  nginx:stable-alpine

# 5. Run parity tests
docker run --rm --network host \
  -e PHP_BASE_URL=http://127.0.0.1:8081 \
  -e GO_BASE_URL=http://127.0.0.1:8082 \
  -e DB_SERVER=127.0.0.1 -e DB_PORT=3308 -e DB_NAME=task-list \
  -e DB_USER=secretuser -e DB_PASSWORD=thisisasupersecretpassworddontyouthink \
  -e PARITY_REPORT_PATH=/app/REPORT.md \
  go-api-test:latest \
  go test -v -count=1 ./tests/parity/...

# 6. Clean up
docker rm -f parity-mysql parity-go-api parity-php-fpm parity-nginx
```

### Parity Report

After running parity tests, a **`REPORT.md`** file is generated containing:

- **Summary table**: total/passed/failed counts
- **Results overview**: one-line per test with status codes and pass/fail
- **Detailed results**: for each test case:
  - The exact HTTP request (method, path, body)
  - The full PHP response (status code, content-type, JSON body)
  - The full Go response (status code, content-type, JSON body)
  - Whether they matched

Set the `PARITY_REPORT_PATH` environment variable to customize where the report is written.

### Parity Matrix

The test cases are defined in [`tests/fixtures/parity_matrix.json`](tests/fixtures/parity_matrix.json). Each entry specifies:

```json
{
  "name": "GET /api/task/1 — happy path",
  "method": "GET",
  "path": "/api/task/1",
  "body": null,
  "expected_status": 200,
  "parity_headers": ["Content-Type"],
  "strict_body_parity": true,
  "ignore_fields": ["created_at", "updated_at"]
}
```

- **`ignore_fields`**: Fields stripped before comparison (timestamps differ between runs)
- **`strict_body_parity`**: When true, JSON bodies must be structurally identical
- **`body_raw`**: For testing invalid JSON (e.g., sending `/` as the body)

### Adding New Parity Tests

1. Add an entry to `tests/fixtures/parity_matrix.json`
2. Add the route to `tests/fixtures/route_inventory.json` (if it's a new endpoint)
3. Run `make test-parity` — the coverage meta-test will fail if any route is missing

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PHP_BASE_URL` | *(required)* | Base URL of the PHP API |
| `GO_BASE_URL` | *(required)* | Base URL of the Go API |
| `DB_SERVER` | `127.0.0.1` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_NAME` | `task-list` | MySQL database name |
| `DB_USER` | `secretuser` | MySQL user |
| `DB_PASSWORD` | `thisisasupersecretpassworddontyouthink` | MySQL password |
| `PARITY_REPORT_PATH` | `REPORT.md` | Path to write the parity report |

## Clean Up

```bash
make clean
```
