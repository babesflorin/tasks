package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetTask mirrors PHP's testGetTask — GET /api/task/1 returns task data.
func TestGetTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	rec := makeRequest(t, router, "GET", "/api/task/1", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	data := body["data"].(map[string]interface{})
	assert.Equal(t, "Task name 0", data["name"])
	assert.Equal(t, "Task description 0", data["description"])
	assert.Equal(t, time.Now().Format("2006-01-02"), data["when"])
	assert.Equal(t, false, data["done"])
}

// TestGetTaskNotFound mirrors PHP's testGetTaskNotFound — GET /api/task/9999999 returns 404.
func TestGetTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	rec := makeRequest(t, router, "GET", "/api/task/9999999", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Empty(t, body["data"])
	assert.Equal(t, "Task not found!", body["error"])
}

// TestAddTask mirrors PHP's testAddTask — POST /api/task creates a new task.
func TestAddTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	taskData := map[string]string{
		"name":        "Task name 6",
		"description": "Task description 6",
		"when":        time.Now().AddDate(0, 0, 5).Format("2006-01-02"),
	}
	reqBody, _ := json.Marshal(taskData)

	rec := makeRequest(t, router, "POST", "/api/task", reqBody)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	data := body["data"].(map[string]interface{})
	assert.Equal(t, taskData["name"], data["name"])
	assert.Equal(t, taskData["description"], data["description"])
	assert.Equal(t, taskData["when"], data["when"])
	assert.Equal(t, false, data["done"])
}

// TestUpdateTask mirrors PHP's testUpdateTask — PUT /api/task updates an existing task.
func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	taskData := map[string]interface{}{
		"id":          2,
		"name":        "Task name 22",
		"description": "Task description 44",
		"when":        time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
	}
	reqBody, _ := json.Marshal(taskData)

	rec := makeRequest(t, router, "PUT", "/api/task", reqBody)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	data := body["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["id"])
	assert.Equal(t, "Task name 22", data["name"])
	assert.Equal(t, "Task description 44", data["description"])
	assert.Equal(t, taskData["when"], data["when"])
	assert.Equal(t, false, data["done"])
}

// TestUpdateTaskNotFound mirrors PHP's testUpdateTaskNotFound — PUT /api/task with non-existent ID returns 404.
func TestUpdateTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	taskData := map[string]interface{}{
		"id":          9999999,
		"name":        "Task name 22",
		"description": "Task description 44",
		"when":        time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
	}
	reqBody, _ := json.Marshal(taskData)

	rec := makeRequest(t, router, "PUT", "/api/task", reqBody)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Empty(t, body["data"])
	assert.Equal(t, "Task not found!", body["error"])
}

// TestUpdateTaskInvalidData mirrors PHP's testUpdateTaskInvalidData — PUT /api/task with empty body returns 400.
func TestUpdateTaskInvalidData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	reqBody := []byte("{}")

	rec := makeRequest(t, router, "PUT", "/api/task", reqBody)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Empty(t, body["data"])
	assert.Equal(t, "Task is not valid!", body["error"])

	expectedMessages := []interface{}{
		"Task name is not valid!",
		"We need an id to know which entity to update!",
		"Task description is not valid!",
		"Task must have a date!",
	}
	assert.Equal(t, expectedMessages, body["messages"])
}

// TestUpdateInvalidJson mirrors PHP's testUpdateInvalidJson — PUT /api/task with invalid JSON returns 400.
func TestUpdateInvalidJson(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	rec := makeRequest(t, router, "PUT", "/api/task", []byte("/"))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Empty(t, body["data"])
	assert.Equal(t, "Request must be json!", body["error"])
}

// TestGetTasks mirrors PHP's testGetTasks — GET /api/task returns all 5 seeded tasks.
func TestGetTasks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	rec := makeRequest(t, router, "GET", "/api/task", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	data := body["data"].([]interface{})
	assert.NotEmpty(t, data)
	assert.Len(t, data, 5)
}

// TestCompleteTask mirrors PHP's testCompleteTask — PUT /api/task/1/complete marks task as done.
func TestCompleteTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	rec := makeRequest(t, router, "PUT", "/api/task/1/complete", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	data := body["data"].(map[string]interface{})
	assert.Equal(t, "Task name 0", data["name"])
	assert.Equal(t, "Task description 0", data["description"])
	assert.Equal(t, time.Now().Format("2006-01-02"), data["when"])
	assert.Equal(t, true, data["done"])
}

// TestCompleteTaskNotFound mirrors PHP's testCompleteTaskNotFound — PUT /api/task/9999999/complete returns 404.
func TestCompleteTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	rec := makeRequest(t, router, "PUT", "/api/task/9999999/complete", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Empty(t, body["data"])
	assert.Equal(t, "Task not found!", body["error"])
}

// TestDeleteTask mirrors PHP's testDeleteTask — DELETE /api/task/5 returns the deleted task.
func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	rec := makeRequest(t, router, "DELETE", "/api/task/5", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	data := body["data"].(map[string]interface{})
	assert.Equal(t, float64(5), data["id"])
	assert.Equal(t, "Task name 4", data["name"])
	assert.Equal(t, "Task description 4", data["description"])
	assert.Equal(t, time.Now().AddDate(0, 0, 4).Format("2006-01-02"), data["when"])
	assert.Equal(t, false, data["done"])
}

// TestDeleteTaskNotFound mirrors PHP's testDeleteTaskNotFound — DELETE /api/task/9999 returns 404.
func TestDeleteTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateAll(t, db)
	seedFixtures(t, db)
	router := setupRouter(t, db)

	rec := makeRequest(t, router, "DELETE", fmt.Sprintf("/api/task/%d", 9999), nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Empty(t, body["data"])
	assert.Equal(t, "Task not found!", body["error"])
}
