package e2e

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/workspace/go-api/internal/handler"
	"github.com/workspace/go-api/internal/repository"
	"github.com/workspace/go-api/internal/service"
)

// setupTestDB connects to the MySQL database using environment variables
// and runs the migration to create the task table.
func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	host := getEnvOrDefault("DB_SERVER", "mysql")
	port := getEnvOrDefault("DB_PORT", "3306")
	name := getEnvOrDefault("DB_NAME", "task-list-test")
	user := getEnvOrDefault("DB_USER", "secretuser")
	pass := getEnvOrDefault("DB_PASSWORD", "thisisasupersecretpassworddontyouthink")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, pass, host, port, name)

	db, err := sqlx.Connect("mysql", dsn)
	require.NoError(t, err, "failed to connect to test database")

	// Run migration
	migration := `CREATE TABLE IF NOT EXISTS task (
		id INT AUTO_INCREMENT NOT NULL,
		name VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL,
		` + "`when`" + ` DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at DATETIME DEFAULT NULL,
		done TINYINT(1) DEFAULT '0' NOT NULL,
		PRIMARY KEY(id)
	) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ENGINE = InnoDB`

	_, err = db.Exec(migration)
	require.NoError(t, err, "failed to run migration")

	return db
}

// setupRouter wires up the full handler stack using a real database connection.
// Returns an http.Handler ready for httptest.
func setupRouter(t *testing.T, db *sqlx.DB) http.Handler {
	t.Helper()

	repo := repository.NewTaskRepository(db)
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	return r
}

// seedFixtures inserts 5 test tasks matching PHP's TasksFixtures.php exactly:
// "Task name 0" through "Task name 4", dates today+0 through today+4.
func seedFixtures(t *testing.T, db *sqlx.DB) {
	t.Helper()

	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		when := time.Now().AddDate(0, 0, i)
		_, err := db.Exec(
			"INSERT INTO task (name, description, `when`, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			fmt.Sprintf("Task name %d", i),
			fmt.Sprintf("Task description %d", i),
			when,
			false,
			now,
			now,
		)
		require.NoError(t, err, "failed to seed fixture %d", i)
	}
}

// truncateAll removes all tasks and resets the auto-increment counter.
func truncateAll(t *testing.T, db *sqlx.DB) {
	t.Helper()
	_, err := db.Exec("TRUNCATE TABLE task")
	require.NoError(t, err, "failed to truncate task table")
}

// makeRequest builds and executes an HTTP request against the router using httptest.
func makeRequest(t *testing.T, router http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
