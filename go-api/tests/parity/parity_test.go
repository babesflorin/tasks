package parity

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParity_AllEndpoints loads the parity matrix and sends identical requests
// to both the PHP and Go APIs, comparing their responses.
func TestParity_AllEndpoints(t *testing.T) {
	phpURL := os.Getenv("PHP_BASE_URL")
	goURL := os.Getenv("GO_BASE_URL")

	if phpURL == "" || goURL == "" {
		t.Skip("PHP_BASE_URL and GO_BASE_URL must be set for parity tests")
	}

	db := connectDB(t)
	defer db.Close()

	matrix := loadParityMatrix(t)

	for _, tc := range matrix {
		t.Run(tc.Name, func(t *testing.T) {
			// Seed shared DB before each test
			truncateAll(t, db)
			seedFixtures(t, db)

			// Build identical request body
			body, err := buildRequestBody(tc)
			require.NoError(t, err)

			// Execute against PHP
			phpResp, phpBody := executeHTTPRequest(t, phpURL, tc.Method, tc.Path, body)
			// Execute against Go
			goResp, goBody := executeHTTPRequest(t, goURL, tc.Method, tc.Path, body)

			// 1. Status code parity
			assert.Equal(t, phpResp.StatusCode, goResp.StatusCode,
				"status code mismatch for %s %s\nPHP body: %s\nGo body: %s",
				tc.Method, tc.Path, phpBody, goBody)

			// 2. Header parity
			for _, h := range tc.ParityHeaders {
				phpHeader := phpResp.Header.Get(h)
				goHeader := goResp.Header.Get(h)
				// Normalize: both should contain "application/json"
				assert.True(t,
					strings.Contains(phpHeader, "json") == strings.Contains(goHeader, "json"),
					"header %s mismatch: PHP=%q Go=%q", h, phpHeader, goHeader)
			}

			// 3. Body parity
			if tc.StrictBodyParity {
				normalizedPHP := normalizeJSON(t, phpBody, tc.IgnoreFields)
				normalizedGo := normalizeJSON(t, goBody, tc.IgnoreFields)
				assert.JSONEq(t, normalizedPHP, normalizedGo,
					"body mismatch for %s %s\nPHP (normalized): %s\nGo  (normalized): %s",
					tc.Method, tc.Path, normalizedPHP, normalizedGo)
			}
		})
	}
}

// TestParity_CoverageCompleteness verifies that every route in the route inventory
// has at least one parity test case in the parity matrix.
func TestParity_CoverageCompleteness(t *testing.T) {
	matrix := loadParityMatrix(t)
	routes := loadRouteInventory(t)

	// Build set of covered routes from the parity matrix.
	// Normalize paths: replace actual IDs with {taskId} pattern.
	coveredRoutes := map[string]bool{}
	for _, tc := range matrix {
		key := tc.Method + " " + normalizePath(tc.Path)
		coveredRoutes[key] = true
	}

	for _, r := range routes {
		key := r.Method + " " + r.Path
		assert.True(t, coveredRoutes[key],
			"route %s has no parity test", key)
	}
}

// normalizePath replaces numeric path segments with {taskId} for route matching.
func normalizePath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if isNumeric(part) {
			parts[i] = "{taskId}"
		}
	}
	return strings.Join(parts, "/")
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
