# Task List E2E Tests

This directory contains comprehensive E2E (End-to-End) tests for the Task Management REST API. These tests are designed to verify 100% coverage of the PHP API functionality to assist in Go migration.

## Overview

The E2E tests are written in Go using the standard `testing` package and `testify` for assertions. They test all REST API endpoints by making HTTP requests to the running service.

## Architecture

- **task_test.go** - Main test file with all test cases
- **cmd/server.go** - Mock server implementing the same API spec for testing
- **Dockerfile** - Container configuration for running tests in Docker

## Running Tests

### Prerequisites

- Go 1.21+
- Running API server (or use the included mock server)

### Run Tests Locally

```bash
# Start the mock server (in one terminal)
go run cmd/server.go

# Run tests (in another terminal)
go test -v ./...
```

### Run Tests in Docker

```bash
docker build -t e2e-tests ./e2e
docker run e2e-tests
```

### Run with Docker Compose

```bash
# From project root
docker compose up -d
docker compose run e2e-tests
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BASE_URL` | Base URL of the API server | `http://localhost:8080/api/task` |
| `SKIP_E2E` | Skip E2E tests when set to `true` | - |

## Test Coverage

The tests cover 100% of the PHP API functionality:

### Endpoints Covered

| Method | Endpoint | Description | Test Functions |
|--------|----------|-------------|----------------|
| POST | `/api/task` | Create a new task | `TestAddTask` |
| GET | `/api/task` | Get all tasks | `TestGetTasks` |
| GET | `/api/task/{id}` | Get a specific task | `TestGetTask` |
| PUT | `/api/task` | Update an existing task | `TestUpdateTask` |
| PUT | `/api/task/{id}/complete` | Mark task as completed | `TestCompleteTask` |
| DELETE | `/api/task/{id}` | Delete a task | `TestDeleteTask` |

### Test Scenarios

#### TestAddTask
- ✅ Success case - create task with valid data
- ✅ Invalid JSON body
- ✅ Missing name validation
- ✅ Missing description validation  
- ✅ Missing when date validation

#### TestGetTasks
- ✅ Get all tasks
- ✅ Filter by completion status (areDone)
- ✅ Filter by date (when)

#### TestGetTask
- ✅ Success case - get existing task
- ✅ Not found - task doesn't exist

#### TestUpdateTask
- ✅ Success case - update task with valid data

#### TestCompleteTask
- ✅ Success case - mark task as completed
- ✅ Not found - task doesn't exist

#### TestDeleteTask
- ✅ Success case - delete existing task
- ✅ Not found - task doesn't exist

## API Specification

### Request/Response Formats

#### TaskRequest (POST /api/task)
```json
{
    "name": "Task Name",
    "description": "Task Description", 
    "when": "2024-12-31"
}
```

#### TaskUpdateRequest (PUT /api/task)
```json
{
    "id": 1,
    "name": "Updated Name",
    "description": "Updated Description",
    "when": "2024-12-31"
}
```

#### TaskResponse
```json
{
    "data": {
        "id": 1,
        "name": "Task Name",
        "description": "Task Description",
        "when": "2024-12-31",
        "done": false,
        "createdAt": "2024-01-01T10:00:00Z",
        "updatedAt": "2024-01-01T10:00:00Z"
    }
}
```

### Query Parameters (GET /api/task)

| Parameter | Type | Description |
|-----------|------|-------------|
| `areDone` | boolean | Filter by completion status |
| `when` | string (Y-m-d) | Filter by due date |

### Error Responses

- **400 Bad Request**: Validation errors
- **404 Not Found**: Task not found

## Migration Notes

These E2E tests are specifically designed to:

1. **Verify Go Migration** - Run against the Go implementation to ensure it matches PHP behavior
2. **Test-Driven Development** - Use tests to drive the Go implementation
3. **Regression Testing** - Ensure no functionality is lost during migration
4. **API Contract** - Document the expected API behavior

## Best Practices

- Tests use dynamic base URL from environment variable
- Tests handle both successful and error responses
- Tests use flexible type assertions to handle JSON number variations
- Tests include proper cleanup (no shared state between tests)