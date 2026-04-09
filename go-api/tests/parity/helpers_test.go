package parity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

// ParityTestCase represents one entry in the parity_matrix.json.
type ParityTestCase struct {
	Name              string                 `json:"name"`
	Method            string                 `json:"method"`
	Path              string                 `json:"path"`
	Body              interface{}            `json:"body"`
	BodyRaw           string                 `json:"body_raw"`
	ExpectedStatus    int                    `json:"expected_status"`
	ParityHeaders     []string               `json:"parity_headers"`
	StrictBodyParity  bool                   `json:"strict_body_parity"`
	IgnoreFields      []string               `json:"ignore_fields"`
}

// RouteEntry represents a route in the route_inventory.json.
type RouteEntry struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// loadParityMatrix reads and parses the parity matrix JSON file.
func loadParityMatrix(t *testing.T) []ParityTestCase {
	t.Helper()
	data, err := os.ReadFile("tests/fixtures/parity_matrix.json")
	if err != nil {
		// Try relative to /app in docker
		data, err = os.ReadFile("/app/tests/fixtures/parity_matrix.json")
		require.NoError(t, err, "failed to read parity_matrix.json")
	}

	var matrix []ParityTestCase
	require.NoError(t, json.Unmarshal(data, &matrix))
	return matrix
}

// loadRouteInventory reads and parses the route inventory JSON file.
func loadRouteInventory(t *testing.T) []RouteEntry {
	t.Helper()
	data, err := os.ReadFile("tests/fixtures/route_inventory.json")
	if err != nil {
		data, err = os.ReadFile("/app/tests/fixtures/route_inventory.json")
		require.NoError(t, err, "failed to read route_inventory.json")
	}

	var routes []RouteEntry
	require.NoError(t, json.Unmarshal(data, &routes))
	return routes
}

// connectDB connects to the shared MySQL database.
func connectDB(t *testing.T) *sqlx.DB {
	t.Helper()

	host := getEnvOrDefault("DB_SERVER", "127.0.0.1")
	port := getEnvOrDefault("DB_PORT", "3306")
	name := getEnvOrDefault("DB_NAME", "task-list")
	user := getEnvOrDefault("DB_USER", "secretuser")
	pass := getEnvOrDefault("DB_PASSWORD", "thisisasupersecretpassworddontyouthink")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, pass, host, port, name)

	db, err := sqlx.Connect("mysql", dsn)
	require.NoError(t, err, "failed to connect to shared test database")
	return db
}

// seedFixtures inserts 5 test tasks matching PHP's TasksFixtures.php.
func seedFixtures(t *testing.T, db *sqlx.DB) {
	t.Helper()
	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		when := time.Now().AddDate(0, 0, i)
		_, err := db.Exec(
			"INSERT INTO task (name, description, `when`, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			fmt.Sprintf("Task name %d", i),
			fmt.Sprintf("Task description %d", i),
			when, false, now, now,
		)
		require.NoError(t, err, "failed to seed fixture %d", i)
	}
}

// truncateAll removes all tasks and resets auto-increment.
func truncateAll(t *testing.T, db *sqlx.DB) {
	t.Helper()
	_, err := db.Exec("TRUNCATE TABLE task")
	require.NoError(t, err)
}

// buildRequestBody prepares the request body for a parity test case.
// It handles dynamic date placeholders like FUTURE_DATE_5.
func buildRequestBody(tc ParityTestCase) (string, error) {
	if tc.BodyRaw != "" {
		return tc.BodyRaw, nil
	}
	if tc.Body == nil {
		return "", nil
	}

	bodyBytes, err := json.Marshal(tc.Body)
	if err != nil {
		return "", err
	}

	body := string(bodyBytes)

	// Replace date placeholders with actual future dates
	for i := 1; i <= 10; i++ {
		placeholder := fmt.Sprintf("FUTURE_DATE_%d", i)
		futureDate := time.Now().AddDate(0, 0, i).Format("2006-01-02")
		body = strings.ReplaceAll(body, placeholder, futureDate)
	}

	return body, nil
}

// executeHTTPRequest sends an HTTP request and returns the response.
func executeHTTPRequest(t *testing.T, baseURL, method, path, body string) (*http.Response, string) {
	t.Helper()

	url := baseURL + path
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	require.NoError(t, err)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

// normalizeJSON removes specified fields from a JSON string for comparison.
func normalizeJSON(t *testing.T, jsonStr string, ignoreFields []string) string {
	t.Helper()

	if len(ignoreFields) == 0 {
		return jsonStr
	}

	var obj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return jsonStr
	}

	removeFields(obj, ignoreFields)

	result, err := json.Marshal(obj)
	require.NoError(t, err)
	return string(result)
}

// removeFields recursively removes specified fields from a JSON structure.
func removeFields(obj interface{}, fields []string) {
	switch v := obj.(type) {
	case map[string]interface{}:
		for _, f := range fields {
			delete(v, f)
		}
		for _, val := range v {
			removeFields(val, fields)
		}
	case []interface{}:
		for _, item := range v {
			removeFields(item, fields)
		}
	}
}

// jsonEqual compares two JSON strings for structural equality (order-insensitive).
func jsonEqual(a, b string) bool {
	var objA, objB interface{}
	if err := json.Unmarshal([]byte(a), &objA); err != nil {
		return a == b
	}
	if err := json.Unmarshal([]byte(b), &objB); err != nil {
		return a == b
	}
	normA, _ := json.Marshal(objA)
	normB, _ := json.Marshal(objB)
	return string(normA) == string(normB)
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
