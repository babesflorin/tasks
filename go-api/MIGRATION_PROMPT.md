# PHP → Go API Migration Prompt

> Reusable prompt for migrating any PHP API application to Go with parity testing.
> Paste this prompt along with a pointer to your PHP codebase (or include an AGENTS.md describing the project structure).

---

You are migrating a PHP API application to Go. Follow this exact process:

## Phase 1 — Discovery & Inventory

Before writing any code, scan the PHP codebase and produce:

### 1.1 Route Map

A table of every endpoint:

| # | Method | Path | Controller::method | Request Body | Query Params | Response Format | Status Codes | Side Effects |
|---|--------|------|--------------------|-------------|-------------|----------------|-------------|-------------|

### 1.2 Response Format Analysis

For each endpoint, document the exact JSON structure returned on success and every error case. Pay attention to:
- Wrapper structures (e.g., Fractal's `{"data": ...}`)
- Timestamp serialization format (PHP's `json_encode(\DateTime)` produces `{"date":"...","timezone_type":3,"timezone":"UTC"}`)
- Error response format (field names, which fields appear in which error types)
- Null vs empty string vs absent fields
- HTTP status codes that don't follow REST conventions (e.g., POST returning 200 instead of 201)

### 1.3 Validation Rules

Document every validation rule with its exact error message string, including typos. The Go implementation must produce byte-identical error messages.

### 1.4 Bug Inventory

Document any bugs found in the PHP code. Policy: fix bugs in PHP first, implement Go correctly from the start. Both codebases stay in sync.

### 1.5 Middleware & External Dependencies

List every middleware (auth, CORS, rate-limit, logging, etc.), database schema, message queues, caches, external API calls, cron jobs.

### 1.6 Docker Infrastructure

Document the full docker-compose.yml, all Dockerfiles, volumes, networks, health-checks, environment variables.

---

## Phase 2 — Fix PHP Bugs

Fix any bugs found during discovery. Run existing PHP tests to verify fixes don't break anything. Commit and push.

---

## Phase 3 — Go Project Scaffolding

Create the Go project with this structure:

```
go-api/
├── cmd/server/main.go                 # Entry point, DI wiring
├── internal/
│   ├── config/                        # Env/config loading
│   ├── handler/                       # HTTP handlers (1:1 with PHP controllers)
│   │   ├── response.go               # Request/Response DTOs matching PHP format exactly
│   │   └── <resource>_handler.go     # Route handlers
│   ├── model/                         # Domain models matching PHP entities
│   ├── repository/                    # DB access (sqlx)
│   ├── service/                       # Business logic + validation
│   └── database/                      # DB connection + migration runner
├── migrations/                        # Same schema as PHP (SQL files)
├── docker/
│   ├── Dockerfile                     # Multi-stage build
│   └── Dockerfile.test               # Test image
├── docker-compose.yml                 # Standalone: Go API + backing services
├── docker-compose.e2e.yml            # Go E2E tests + backing services
├── docker-compose.test.yml           # Parity: PHP + Go + shared backing services
├── Makefile
├── go.mod
├── REPORT.md                          # Auto-generated parity report
└── tests/
    ├── e2e/                           # Go E2E tests (1:1 with PHP tests)
    ├── parity/                        # Cross-parity tests (PHP vs Go)
    └── fixtures/                      # route_inventory.json, parity_matrix.json
```

**Tech stack**: chi (router), sqlx + appropriate DB driver, testify (assertions). Match the PHP app's database (MySQL/PostgreSQL).

---

## Phase 4 — Go Core Implementation

Implement each layer, matching PHP behavior exactly:

1. **Response structs** — Must produce byte-identical JSON to PHP. Replicate PHP's serialization quirks (DateTime format, Fractal wrapping, error response structure).
2. **Validator** — Same validation rules, same error messages (including typos), same ordering of error messages in arrays.
3. **Repository** — Same SQL queries, same filtering behavior.
4. **Service** — Same business logic, same error types mapping to same HTTP status codes.
5. **Handlers** — Same routes, same request parsing, same error handling (replicate PHP's ExceptionListener behavior).
6. **Entry point** — Include DB connection retry loop (replaces PHP's wait-for-it.sh pattern).

Verify: `go build ./cmd/server` succeeds.

---

## Phase 5 — Docker Configuration

The Go docker-compose MUST:
- Use the SAME backing service images and versions as PHP
- Use the SAME environment variable names
- Include health-checks on all services
- Use host networking if Docker bridge networking is unavailable

Create three docker-compose files:
- `docker-compose.yml` — Standalone Go API
- `docker-compose.e2e.yml` — Go E2E tests against real database
- `docker-compose.test.yml` — Parity mode: both PHP and Go APIs running simultaneously sharing the same database

---

## Phase 6 — Three Test Suites

### 6.1 Go E2E Tests

1:1 mapping with every PHP functional test. Each test:
- Seeds fixtures (matching PHP's test fixtures exactly)
- Sends HTTP request via `net/http/httptest`
- Asserts status code, Content-Type, and specific JSON field values
- Truncates DB after each test

### 6.2 Parity Tests

For each endpoint, send identical requests to both PHP and Go APIs and compare:
1. Status codes must match
2. Content-Type headers must match
3. Response bodies must be structurally identical (after normalizing timestamps/auto-increment IDs)

Key rules:
- Seed DB fresh before EACH API request (so side effects from one don't affect the other)
- Include a coverage completeness meta-test that verifies every route has at least one parity test

### 6.3 Parity Report

Auto-generate a `REPORT.md` after parity tests containing:
- Summary table (total/passed/failed)
- Results overview (one-line per test)
- Detailed results: for each test, the exact request sent, the full PHP response, and the full Go response

### 6.4 Test Fixtures

- `tests/fixtures/route_inventory.json` — All routes
- `tests/fixtures/parity_matrix.json` — All parity test cases with fields:

```json
{
  "name": "GET /api/resource/1 — happy path",
  "method": "GET",
  "path": "/api/resource/1",
  "body": null,
  "expected_status": 200,
  "parity_headers": ["Content-Type"],
  "strict_body_parity": true,
  "ignore_fields": ["created_at", "updated_at"]
}
```

---

## Phase 7 — Makefile & CI

```makefile
test-e2e:     # Go E2E tests in Docker against real DB
test-php:     # PHP tests in Docker
test-parity:  # Both APIs + shared DB, generates REPORT.md
test-all:     # All three suites
clean:        # Tear down all Docker resources
```

GitHub Actions workflow running all three test suites.

---

## Normalization Rules for Parity

When comparing PHP vs Go responses, normalize:
- **Timestamps**: Strip from comparison (format matches but values differ between runs)
- **Auto-increment IDs**: Ignore or match by position
- **Key ordering**: Use order-insensitive JSON comparison
- **Null vs absent**: Document convention and enforce in both

---

## Acceptance Criteria

- [ ] PHP bugs fixed, PHP tests pass after fix
- [ ] `go build ./cmd/server` succeeds
- [ ] Go E2E tests pass in Docker against real database
- [ ] Every PHP endpoint has a Go equivalent at the same path
- [ ] Response JSON structure matches PHP exactly
- [ ] Error messages are byte-identical to PHP
- [ ] Validation rules produce identical error arrays in same order
- [ ] Docker configs mirror PHP's infrastructure
- [ ] Go API runs in Docker with same backing services as PHP
- [ ] Parity tests pass (both APIs produce identical responses)
- [ ] REPORT.md generated with full request/response comparison
- [ ] Bugs found are fixed in both codebases
