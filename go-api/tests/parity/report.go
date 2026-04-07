package parity

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ReportEntry captures the full request/response details for a single parity test case.
type ReportEntry struct {
	TestName       string `json:"test_name"`
	Method         string `json:"method"`
	Path           string `json:"path"`
	RequestBody    string `json:"request_body,omitempty"`
	PHPStatus      int    `json:"php_status"`
	PHPBody        string `json:"php_body"`
	PHPContentType string `json:"php_content_type"`
	GoStatus       int    `json:"go_status"`
	GoBody         string `json:"go_body"`
	GoContentType  string `json:"go_content_type"`
	Match          bool   `json:"match"`
	Notes          string `json:"notes,omitempty"`
}

// Report collects all parity test results and writes them to REPORT.md.
type Report struct {
	mu      sync.Mutex
	entries []ReportEntry
	startAt time.Time
}

// NewReport creates a new empty report.
func NewReport() *Report {
	return &Report{
		entries: []ReportEntry{},
		startAt: time.Now(),
	}
}

// Add appends a test result entry to the report.
func (r *Report) Add(entry ReportEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, entry)
}

// WriteMarkdown generates REPORT.md from all collected entries.
func (r *Report) WriteMarkdown(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var sb strings.Builder

	sb.WriteString("# Parity Test Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated**: %s\n\n", r.startAt.UTC().Format("2006-01-02 15:04:05 UTC")))
	sb.WriteString(fmt.Sprintf("**Duration**: %s\n\n", time.Since(r.startAt).Round(time.Millisecond)))

	// Summary table
	passed := 0
	failed := 0
	for _, e := range r.entries {
		if e.Match {
			passed++
		} else {
			failed++
		}
	}
	total := passed + failed

	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("| Total | Passed | Failed |\n"))
	sb.WriteString(fmt.Sprintf("|-------|--------|--------|\n"))
	sb.WriteString(fmt.Sprintf("| %d | %d | %d |\n\n", total, passed, failed))

	// Overview table
	sb.WriteString("## Results Overview\n\n")
	sb.WriteString("| # | Test | Method | Path | PHP Status | Go Status | Match |\n")
	sb.WriteString("|---|------|--------|------|------------|-----------|-------|\n")
	for i, e := range r.entries {
		matchIcon := "✅"
		if !e.Match {
			matchIcon = "❌"
		}
		sb.WriteString(fmt.Sprintf("| %d | %s | `%s` | `%s` | %d | %d | %s |\n",
			i+1, e.TestName, e.Method, e.Path, e.PHPStatus, e.GoStatus, matchIcon))
	}
	sb.WriteString("\n")

	// Detailed results
	sb.WriteString("## Detailed Results\n\n")
	for i, e := range r.entries {
		matchIcon := "✅ PASS"
		if !e.Match {
			matchIcon = "❌ FAIL"
		}
		sb.WriteString(fmt.Sprintf("### %d. %s — %s\n\n", i+1, e.TestName, matchIcon))
		sb.WriteString(fmt.Sprintf("**Request**: `%s %s`\n\n", e.Method, e.Path))

		if e.RequestBody != "" {
			sb.WriteString("**Request Body**:\n")
			sb.WriteString("```json\n")
			sb.WriteString(prettyJSON(e.RequestBody))
			sb.WriteString("\n```\n\n")
		}

		// PHP response
		sb.WriteString(fmt.Sprintf("**PHP Response** (status %d, `%s`):\n", e.PHPStatus, e.PHPContentType))
		sb.WriteString("```json\n")
		sb.WriteString(prettyJSON(e.PHPBody))
		sb.WriteString("\n```\n\n")

		// Go response
		sb.WriteString(fmt.Sprintf("**Go Response** (status %d, `%s`):\n", e.GoStatus, e.GoContentType))
		sb.WriteString("```json\n")
		sb.WriteString(prettyJSON(e.GoBody))
		sb.WriteString("\n```\n\n")

		if e.Notes != "" {
			sb.WriteString(fmt.Sprintf("**Notes**: %s\n\n", e.Notes))
		}

		sb.WriteString("---\n\n")
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

// prettyJSON attempts to pretty-print a JSON string. Falls back to raw string.
func prettyJSON(s string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return s
	}
	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return s
	}
	return string(pretty)
}
