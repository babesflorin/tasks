# AGENTS.md - Symfony Task Management REST API

## 1. Project Overview

This is a **Symfony 5.0 REST API** project for a **Task Management System**. The project follows a layered architecture with clear separation between Domain, Presentation, and Infrastructure layers.

### Tech Stack
- **PHP**: 7.2+
- **Framework**: Symfony 5.0
- **ORM**: Doctrine ORM
- **API Documentation**: NelmioApiDocBundle (Swagger UI at `/api/doc`)
- **Data Transformation**: League/Fractal
- **Testing**: PHPUnit with custom TestCase
- **Code Style**: PHP CodeSniffer
- **Container**: Docker

### Key Dependencies (from composer.json)
```json
{
    "doctrine/annotations": "^1.8",
    "league/fractal": "^0.19.2",
    "nelmio/api-doc-bundle": "^3.6",
    "sensio/framework-extra-bundle": "^5.5",
    "symfony/framework-bundle": "5.0.*",
    "symfony/orm-pack": "^1.0"
}
```

---

## 2. Project Structure

```
/workspace/project/tasks/
├── bin/                          # Console commands
│   └── console
├── config/                       # Configuration files
│   ├── packages/                 # Environment-specific config
│   ├── bundles.php               # Bundle registration
│   ├── routes.yaml               # Route definitions
│   ├── services.yaml             # Service container configuration
│   └── bootstrap.php
├── src/
│   ├── Domain/                   # Business logic layer
│   │   ├── Dto/                  # Data Transfer Objects (TaskDto, TaskCollectionDto)
│   │   ├── Entity/               # Domain entities (Task, TaskCollection)
│   │   ├── Exception/            # Custom exceptions
│   │   ├── Repository/           # Repository interfaces
│   │   ├── Service/              # Domain services (TaskService)
│   │   ├── Transformer/          # DTO to Entity transformers
│   │   └── Validator/            # Domain validators (TaskValidator)
│   ├── Infrastructure/           # External concerns
│   │   ├── Mapping/              # Doctrine mappings
│   │   ├── Migrations/           # Database migrations
│   │   └── Repository/           # Repository implementations
│   ├── Presentation/             # Input/output layer
│   │   ├── Controller/           # REST controllers (TaskController)
│   │   ├── EventListener/        # Event listeners (ExceptionListener)
│   │   ├── Request/              # Request handlers (ParamConverter)
│   │   └── Transformer/          # Response transformers (TaskTransformer)
│   └── Kernel.php                # Symfony kernel
├── tests/
│   ├── DataFixtures/             # Test data fixtures
│   ├── functional/               # Functional tests
│   ├── unit/                     # Unit tests
│   └── bootstrap.php             # Test bootstrap
├── public/                       # Web root
│   └── index.php                 # Entry point
├── docker/                       # Docker configuration
│   ├── build/
│   ├── nginx/
│   ├── php-fpm/
│   └── installapp/
├── hooks/                        # Git hooks
│   ├── pre-commit                # Pre-commit hook (code style checks)
│   └── setup.sh
├── vendor/                       # Composer dependencies
├── .env                          # Environment variables
├── composer.json                 # PHP dependencies
├── phpunit.xml.dist              # PHPUnit configuration
├── phpcs.xml.dist                # CodeSniffer configuration
└── docker-compose.yaml           # Docker services
```

---

## 3. API Endpoints

All endpoints are prefixed with `/api/task` and accept/return JSON.

### Base URL
```
/api/task
```

### Endpoints

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/api/task` | Create a new task | TaskRequest |
| GET | `/api/task` | Get all tasks (filterable) | - |
| GET | `/api/task/{taskId}` | Get a specific task | - |
| PUT | `/api/task` | Update an existing task | TaskUpdateRequest |
| PUT | `/api/task/{taskId}/complete` | Mark a task as completed | - |
| DELETE | `/api/task/{taskId}` | Delete a task | - |

### Query Parameters (GET /api/task)

| Parameter | Type | Description |
|-----------|------|-------------|
| `areDone` | boolean | Filter by completion status |
| `when` | string (Y-m-d) | Filter by due date |

### Request/Response Formats

#### TaskRequest (POST)
```json
{
    "name": "Task Name",
    "description": "Task Description",
    "when": "2024-12-31"
}
```

#### TaskUpdateRequest (PUT)
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
    "id": 1,
    "name": "Task Name",
    "description": "Task Description",
    "when": "2024-12-31",
    "done": false,
    "createdAt": "2024-01-01T10:00:00+00:00",
    "updatedAt": "2024-01-01T10:00:00+00:00"
}
```

#### MultipleTaskResponse
```json
{
    "data": [
        { ... TaskResponse ... },
        { ... TaskResponse ... }
    ]
}
```

### Error Responses

- **400 Bad Request**: Validation errors
- **404 Not Found**: Task not found

---

## 4. Key Services and Dependencies

### TaskService (`src/Domain/Service/TaskService.php`)
The main domain service handling all task operations:

```php
public function addTask(TaskDto $taskDto): TaskDto
public function getAllTasks(?bool $areDone = null, ?\DateTime $when = null): TaskCollectionDto
public function completeTask(int $taskId): TaskDto
public function updateTask(TaskDto $taskDto): TaskDto
public function getTask(int $taskId): TaskDto
public function deleteTask(int $taskId): TaskDto
```

### TaskRepositoryInterface (`src/Domain/Repository/TaskRepositoryInterface.php`)
Defines the contract for data access:

```php
public function saveTask(Task $task): Task
public function getTasks(?bool $areDone = null, ?\DateTime $when = null): TaskCollection
public function findTaskById(int $taskId): ?Task
public function deleteTask(Task $task): bool
```

### TaskRepository (`src/Infrastructure/Repository/TaskRepository.php`)
Doctrine implementation of TaskRepositoryInterface.

### TaskValidator (`src/Domain/Validator/TaskValidator.php`)
Validates task DTOs with the following rules:
- `name`: Required, non-empty string
- `description`: Required, non-empty string
- `when`: Required, valid date in format Y-m-d, must be in the future or today

### Fractal Manager
Configured in `config/services.yaml` for transforming entities to JSON responses.

### TaskTransformer (`src/Presentation/Transformer/TaskTransformer.php`)
Transforms Task entity to JSON-serializable array for API responses.

---

## 5. Development Commands

### Install Dependencies
```bash
composer install
```

### Clear Cache
```bash
php bin/console cache:clear
```

### Run Tests
```bash
# All tests
php bin/phpunit

# Unit tests only
php bin/phpunit --testsuite UnitTests

# Functional tests only
php bin/phpunit --testsuite FunctionalTests

# With coverage
php bin/phpunit --coverage-html coverage-result
```

### Code Style Check
```bash
# Check coding standards
vendor/bin/phpcs --standard=phpcs.xml.dist src/

# Auto-fix coding standards (use with caution)
vendor/bin/phpcbf --standard=phpcs.xml.dist src/
```

### Database Migrations
```bash
# Generate migration
php bin/console make:migration

# Run migrations
php bin/console doctrine:migrations:migrate

# Check migration status
php bin/console doctrine:migrations:status
```

### Load Fixtures
```bash
php bin/console doctrine:fixtures:load
```

### Symfony Console
```bash
# List all commands
php bin/console list

# Debug routes
php bin/console debug:router
```

---

## 6. Database Schema

### Task Entity (`src/Domain/Entity/Task.php`)

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Primary key, auto-increment |
| `name` | string (255) | Task name |
| `description` | text | Task description |
| `when` | date | Due date |
| `done` | boolean | Completion status (default: false) |
| `createdAt` | datetime | Creation timestamp |
| `updatedAt` | datetime | Last update timestamp |

The database migration is at: `src/Infrastructure/Migrations/Version20200308095622.php`

### Doctrine Configuration
- Database URL is configured in `.env` file
- Default: `mysql://secretuser:thisisasupersecretpassworddontyouthink@mysql:3306/task-list?serverVersion=5.7`

---

## 7. Validation Rules

TaskValidator (`src/Domain/Validator/TaskValidator.php`) enforces:

| Field | Validation |
|-------|------------|
| `name` | Required, non-empty string |
| `description` | Required, non-empty string |
| `when` | Required, valid Y-m-d format, date must be >= today |

### Validation Errors
Invalid tasks throw `InvalidTaskException` with an array of error messages.

---

## 8. Testing Patterns

### Test Structure
- **Unit Tests**: `tests/unit/`
- **Functional Tests**: `tests/functional/`

### Custom TestCase (`tests/unit/TestCase.php`)
Provides helper method for creating mocks:
```php
$this->getMock(string $className, array $methods = []): MockObject
```

### Unit Test Examples

**TaskServiceTest** (`tests/unit/Domain/Service/TaskServiceTest.php`):
- Tests all TaskService methods
- Uses mocks for repository, validator, and transformers

**TaskValidatorTest** (`tests/unit/Domain/Validator/TaskValidatorTest.php`):
- Tests validation rules

**TaskEntityTest** (`tests/unit/Domain/Entity/TaskTest.php`):
- Tests domain entity behavior

### Functional Test Example

**TaskControllerTest** (`tests/functional/Presentation/Controller/TaskControllerTest.php`):
- Tests HTTP endpoints
- Uses Symfony's test client

### Fixtures
`tests/DataFixtures/TasksFixtures.php` - Loads sample task data for testing.

---

## 9. Docker Environment

### Services (docker-compose.yaml)

| Service | Description | Port |
|---------|-------------|------|
| `mysql` | MySQL 5.7 database | 3306 (configurable) |
| `php-fpm` | PHP-FPM application server | - |
| `webserver` | Nginx web server | 8101 (configurable) |
| `builder` | Build container | - |
| `installapp` | Application installer | - |

### Environment Variables (.env)
```
APP_ENV=dev
WEBSERVER_PORT=8101
DB_NAME=task-list
DB_USER=secretuser
DB_PASSWORD=thisisasupersecretpassworddontyouthink
DB_ROOT_PASSWORD=iamgroot
DB_SERVER=mysql
DB_PORT=3306
DATABASE_URL=mysql://secretuser:thisisasupersecretpassworddontyouthink@mysql:3306/task-list?serverVersion=5.7
```

### Docker Commands
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

---

## 10. Important Notes

### API Documentation
- Swagger UI is available at: `/api/doc`
- Configuration: `config/packages/nelmio_api_doc.yaml`

### ParamConverter
- `TaskDto` is automatically converted from request body using SensioFrameworkExtraBundle's ParamConverter
- See: `src/Presentation/Request/ParamConverter/TaskConverter.php`

### Exception Handling
- Custom exceptions are handled by `src/Presentation/EventListener/ExceptionListener.php`
- Maps domain exceptions to appropriate HTTP responses

### Code Style
- Pre-commit hook runs PHP CodeSniffer to enforce PSR-12 coding standards
- Configuration: `phpcs.xml.dist`

### Fractal Integration
- League/Fractal is used for transforming entities to JSON
- TaskTransformer handles the transformation logic

---

## 11. Working Guidelines

### Adding a New Feature

1. **Understand the domain**: Start with the Domain layer
2. **Create/Update Entity**: If needed, modify entities in `src/Domain/Entity/`
3. **Create/Update DTO**: Add or modify DTOs in `src/Domain/Dto/`
4. **Create/Update Service**: Add business logic in `src/Domain/Service/`
5. **Create/Update Repository**: Implement data access in `src/Infrastructure/Repository/`
6. **Create Transformer**: Add transformation logic if needed
7. **Create Controller**: Add endpoint in `src/Presentation/Controller/`
8. **Update Routes**: If needed, update route configuration
9. **Write Tests**: Add unit and/or functional tests
10. **Run Code Style Check**: Ensure code follows standards

### Bug Fix Process

1. **Reproduce**: Write a test that reproduces the bug
2. **Identify**: Find the source of the bug in the codebase
3. **Fix**: Implement the fix
4. **Verify**: Run the test to verify the fix
5. **Check**: Run all tests and code style checks

### Code Style Standards
- Follow PSR-12 coding standards
- Run `vendor/bin/phpcs --standard=phpcs.xml.dist src/` before committing
- Use type hints where possible
- Document public methods with docblocks

### Testing Guidelines
- Write unit tests for domain logic
- Write functional tests for controllers
- Mock external dependencies
- Aim for high code coverage

---

## 12. Exception Classes

Located in `src/Domain/Exception/`:

| Exception | Use Case |
|-----------|----------|
| `ValidationException` | General validation errors |
| `InvalidTaskException` | Task validation failures |
| `TaskNotFoundException` | Task not found (404) |
| `CouldNotDeleteException` | Delete operation failures |

---

## 13. Request/Response Flow

1. **Request** → `TaskController` endpoint
2. **ParamConverter** → Converts JSON body to `TaskDto`
3. **TaskService** → Processes business logic
4. **TaskValidator** → Validates input
5. **Repository** → Persists/retrieves data
6. **Transformer** → Converts Entity to DTO
7. **Fractal** → Transforms DTO to JSON
8. **Response** → JSON to client

---

## 14. File Summary

### Source Files (21 PHP files in src/)
- Controllers: 1 (TaskController)
- Entities: 3 (Task, TaskCollection, Collection)
- DTOs: 2 (TaskDto, TaskCollectionDto)
- Services: 1 (TaskService)
- Repositories: 2 (Interface + Implementation)
- Validators: 1 (TaskValidator)
- Transformers: 3 (DTO to Entity, DTO Collection, Presentation)
- Event Listeners: 1 (ExceptionListener)
- Param Converters: 1 (TaskConverter)
- Exceptions: 4 (ValidationException, InvalidTaskException, TaskNotFoundException, CouldNotDeleteException)
- Kernel: 1 (Kernel.php)
- Migrations: 1 (Version20200308095622.php)

## Cursor Cloud specific instructions

### Environment overview

This is a Symfony 5.0 REST API (PHP 7.4) with Docker-based services (MySQL 5.7, PHP-FPM, Nginx). Tests use SQLite in-memory and do not require Docker.

### Running tests (no Docker needed)

```bash
cd /workspace && vendor/bin/simple-phpunit
```

Tests use SQLite via `config/packages/test/doctrine.yaml`, so MySQL is not required.

### Running lint

```bash
vendor/bin/phpcs --standard=phpcs.xml.dist src/
```

### Running the application (Docker)

```bash
sudo docker compose up -d
sudo docker compose exec php-fpm php bin/console doctrine:schema:update --force
```

The API is then available at `http://localhost:8101/api/task`. Swagger UI at `http://localhost:8101/api/doc`.

### Gotchas

- The `composer.lock` was generated with Composer 1. Use `composer config --no-plugins allow-plugins.ocramius/package-versions true` and `composer config --no-plugins allow-plugins.symfony/flex true` before `composer install` on Composer 2.
- The `docker/build/Dockerfile` uses `--ignore-platform-reqs --no-scripts` because the `composer` Docker image ships PHP 8.x+ while the project targets PHP 7.4. The runtime `php-fpm` container uses PHP 7.4.
- `doctrine:migrations:migrate` may fail due to a Doctrine Migrations version incompatibility (`Undefined class constant 'VERSIONS'`). Use `doctrine:schema:update --force` instead for initial DB setup.
- No `bin/phpunit` exists; use `vendor/bin/simple-phpunit` which bootstraps PHPUnit 7.5 via the Symfony PHPUnit Bridge.
- Code coverage requires Xdebug; the "No code coverage driver is available" warning is harmless for running tests.