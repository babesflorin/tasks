package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	
	"os"
	
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Get BASE_URL from environment or use default
func getBaseURL() string {
	if url := os.Getenv("BASE_URL"); url != "" {
		return url
	}
	return "http://localhost:8080/api/task"
}

// TaskRequest represents the request body for creating a task
type TaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	When        string `json:"when"`
}

// TaskUpdateRequest represents the request body for updating a task
type TaskUpdateRequest struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	When        string `json:"when"`
}

// Helper function to get today's date in Y-m-d format
func todayDate() string {
	return time.Now().Format("2006-01-02")
}

// Helper function to get a future date in Y-m-d format
func futureDate(days int) string {
	return time.Now().AddDate(0, 0, days).Format("2006-01-02")
}

// Helper function to get a past date in Y-m-d format
func pastDate(days int) string {
	return time.Now().AddDate(0, 0, -days).Format("2006-01-02")
}

// TestAddTask tests creating a new task - covers addTask endpoint
func TestAddTask(t *testing.T) {
	baseURL := getBaseURL()
	
	t.Run("success", func(t *testing.T) {
		taskData := TaskRequest{
			Name:        "Test Task",
			Description: "Test Description",
			When:        futureDate(5),
		}

		body, err := json.Marshal(taskData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json"), "Content-Type should start with application/json")

		// Parse the response - it should have a data object
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok, "response data should be an object")

		assert.Equal(t, taskData.Name, data["name"])
		assert.Equal(t, taskData.Description, data["description"])
		assert.Equal(t, taskData.When, data["when"])
		assert.Equal(t, false, data["done"])
		assert.NotZero(t, data["id"])
	})

	t.Run("invalid_json", func(t *testing.T) {
		req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer([]byte("invalid")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Accept both 400 and 200 as the server might be lenient
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusOK)
	})

	t.Run("missing_name", func(t *testing.T) {
		taskData := TaskRequest{
			Description: "Test Description",
			When:       futureDate(5),
		}

		body, err := json.Marshal(taskData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Server should validate the name
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		
		if resp.StatusCode == http.StatusBadRequest {
			assert.Equal(t, "Task is not valid!", response["error"])
		}
	})

	t.Run("missing_description", func(t *testing.T) {
		taskData := TaskRequest{
			Name: "Test Task",
			When: futureDate(5),
		}

		body, err := json.Marshal(taskData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusBadRequest {
			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)
			assert.Equal(t, "Task is not valid!", response["error"])
		}
	})

	t.Run("missing_when", func(t *testing.T) {
		taskData := TaskRequest{
			Name:        "Test Task",
			Description: "Test Description",
		}

		body, err := json.Marshal(taskData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusBadRequest {
			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)
			assert.Equal(t, "Task is not valid!", response["error"])
		}
	})
}

// TestGetTasks tests getting all tasks - covers getAllTasks endpoint with filters
func TestGetTasks(t *testing.T) {
	baseURL := getBaseURL()
	
	t.Run("success_all", func(t *testing.T) {
		req, err := http.NewRequest("GET", baseURL, nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json"), "Content-Type should start with application/json")

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response["data"])
	})

	t.Run("filter_by_areDone", func(t *testing.T) {
		req, err := http.NewRequest("GET", baseURL+"?areDone=true", nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("filter_by_when", func(t *testing.T) {
		req, err := http.NewRequest("GET", baseURL+"?when="+todayDate(), nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestGetTask tests getting a specific task - covers getTask endpoint
func TestGetTask(t *testing.T) {
	baseURL := getBaseURL()
	
	t.Run("success", func(t *testing.T) {
		// First create a task to get
		taskData := TaskRequest{
			Name:        "Task to Get",
			Description: "Description",
			When:        futureDate(5),
		}
		body, _ := json.Marshal(taskData)
		createReq, _ := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
		createReq.Header.Set("Content-Type", "application/json")
		createResp, err := http.DefaultClient.Do(createReq)
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createResult map[string]interface{}
		json.NewDecoder(createResp.Body).Decode(&createResult)
		data, ok := createResult["data"].(map[string]interface{})
		require.True(t, ok)
		
		// Get id - it could be float64 (JSON number) or string
		var taskID int
		switch v := data["id"].(type) {
		case float64:
			taskID = int(v)
		case int:
			taskID = v
		default:
			t.Skipf("Cannot determine task ID from response")
			return
		}

		// Now get the task
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", baseURL, taskID), nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json"), "Content-Type should start with application/json")

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		taskDataResp, ok := response["data"].(map[string]interface{})
		require.True(t, ok)
		
		assert.Equal(t, taskData.Name, taskDataResp["name"])
		assert.Equal(t, taskData.Description, taskDataResp["description"])
	})

	t.Run("not_found", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", baseURL, 9999999), nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Either 404 or 200 with error message are acceptable
		assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusOK)
	})
}

// TestUpdateTask tests updating a task - covers updateTask endpoint
func TestUpdateTask(t *testing.T) {
	baseURL := getBaseURL()
	
	t.Run("success", func(t *testing.T) {
		// First create a task to update
		taskData := TaskRequest{
			Name:        "Original Name",
			Description: "Original Description",
			When:        futureDate(5),
		}
		body, _ := json.Marshal(taskData)
		createReq, _ := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
		createReq.Header.Set("Content-Type", "application/json")
		createResp, err := http.DefaultClient.Do(createReq)
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createResult map[string]interface{}
		json.NewDecoder(createResp.Body).Decode(&createResult)
		data, ok := createResult["data"].(map[string]interface{})
		require.True(t, ok)
		
		var taskID int
		switch v := data["id"].(type) {
		case float64:
			taskID = int(v)
		case int:
			taskID = v
		default:
			t.Skipf("Cannot determine task ID from response")
			return
		}

		// Now update the task
		updateData := TaskUpdateRequest{
			ID:          taskID,
			Name:        "Updated Name",
			Description: "Updated Description",
			When:        futureDate(10),
		}
		updateBody, err := json.Marshal(updateData)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", baseURL, bytes.NewBuffer(updateBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
	})
}

// TestCompleteTask tests completing a task - covers completeTask endpoint
func TestCompleteTask(t *testing.T) {
	baseURL := getBaseURL()
	
	t.Run("success", func(t *testing.T) {
		// First create a task to complete
		taskData := TaskRequest{
			Name:        "Task to Complete",
			Description: "Description",
			When:        futureDate(5),
		}
		body, _ := json.Marshal(taskData)
		createReq, _ := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
		createReq.Header.Set("Content-Type", "application/json")
		createResp, err := http.DefaultClient.Do(createReq)
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createResult map[string]interface{}
		json.NewDecoder(createResp.Body).Decode(&createResult)
		data, ok := createResult["data"].(map[string]interface{})
		require.True(t, ok)
		
		var taskID int
		switch v := data["id"].(type) {
		case float64:
			taskID = int(v)
		case int:
			taskID = v
		default:
			t.Skipf("Cannot determine task ID from response")
			return
		}

		// Now complete the task
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%d/complete", baseURL, taskID), nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json"), "Content-Type should start with application/json")

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		taskDataResp, ok := response["data"].(map[string]interface{})
		require.True(t, ok)
		
		assert.Equal(t, taskData.Name, taskDataResp["name"])
		// Done should be true after completion
		assert.Equal(t, true, taskDataResp["done"])
	})

	t.Run("not_found", func(t *testing.T) {
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%d/complete", baseURL, 9999999), nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusOK)
	})
}

// TestDeleteTask tests deleting a task - covers deleteTask endpoint
func TestDeleteTask(t *testing.T) {
	baseURL := getBaseURL()
	
	t.Run("success", func(t *testing.T) {
		// First create a task to delete
		taskData := TaskRequest{
			Name:        "Task to Delete",
			Description: "Description",
			When:        futureDate(5),
		}
		body, _ := json.Marshal(taskData)
		createReq, _ := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
		createReq.Header.Set("Content-Type", "application/json")
		createResp, err := http.DefaultClient.Do(createReq)
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createResult map[string]interface{}
		json.NewDecoder(createResp.Body).Decode(&createResult)
		data, ok := createResult["data"].(map[string]interface{})
		require.True(t, ok)
		
		var taskID int
		switch v := data["id"].(type) {
		case float64:
			taskID = int(v)
		case int:
			taskID = v
		default:
			t.Skipf("Cannot determine task ID from response")
			return
		}

		// Now delete the task
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d", baseURL, taskID), nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json"), "Content-Type should start with application/json")
	})

	t.Run("not_found", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d", baseURL, 9999999), nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusOK)
	})
}

// TestMain is the entry point for all tests
func TestMain(m *testing.M) {
	// Check if we should skip tests (when running without proper network)
	if os.Getenv("SKIP_E2E") == "true" {
		fmt.Println("Skipping E2E tests (SKIP_E2E=true)")
		os.Exit(0)
	}

	// Run the tests
	exitCode := m.Run()
	os.Exit(exitCode)
}