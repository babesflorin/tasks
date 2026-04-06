package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Task represents a task in the system
type Task struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	When        string    `json:"when"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// In-memory store
var tasks = make(map[int]*Task)
var nextID = 1

// Response types
type TaskResponse struct {
	Data *Task `json:"data"`
}

type TasksListResponse struct {
	Data []*Task `json:"data"`
}

type ErrorResponse struct {
	Error   string   `json:"error"`
	Messages []string `json:"messages,omitempty"`
}

// Helper to get today's date in Y-m-d format
func todayDate() string {
	return time.Now().Format("2006-01-02")
}

func main() {
	// Initialize with some sample data
	for i := 0; i < 5; i++ {
		task := &Task{
			ID:          nextID,
			Name:        fmt.Sprintf("Task name %d", i),
			Description: fmt.Sprintf("Task description %d", i),
			When:        todayDate(),
			Done:        false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		tasks[nextID] = task
		nextID++
	}

	// Set up routes
	http.HandleFunc("/api/task", handleTasks)
	http.HandleFunc("/api/task/", handleTaskByID)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		createTask(w, r)
	case "GET":
		getTasks(w, r)
	case "PUT":
		updateTask(w, r)
	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract ID from path
	var taskID int
	fmt.Sscanf(r.URL.Path, "/api/task/%d", &taskID)

	if r.URL.Path == fmt.Sprintf("/api/task/%d/complete", taskID) {
		completeTask(w, r, taskID)
		return
	}

	switch r.Method {
	case "GET":
		getTask(w, r, taskID)
	case "DELETE":
		deleteTask(w, r, taskID)
	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var taskReq struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		When        string `json:"when"`
	}

	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Request must be json!"})
		return
	}

	// Validation
	var messages []string
	if taskReq.Name == "" {
		messages = append(messages, "Task name is not valid!")
	}
	if taskReq.Description == "" {
		messages = append(messages, "Task description is not valid!")
	}
	if taskReq.When == "" {
		messages = append(messages, "Task must have a date!")
	} else if _, err := time.Parse("2006-01-02", taskReq.When); err != nil {
		messages = append(messages, "`when` is not a valid date!")
	} else if taskReq.When < todayDate() {
		messages = append(messages, "You can't do a task in a the past. Or you can?")
	}

	if len(messages) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Task is not valid!", Messages: messages})
		return
	}

	task := &Task{
		ID:          nextID,
		Name:        taskReq.Name,
		Description: taskReq.Description,
		When:        taskReq.When,
		Done:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks[nextID] = task
	nextID++

	json.NewEncoder(w).Encode(TaskResponse{Data: task})
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	areDone := r.URL.Query().Get("areDone")
	when := r.URL.Query().Get("when")

	// Validate 'when' parameter if provided
	if when != "" {
		if _, err := time.Parse("2006-01-02", when); err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "When is not valid!"})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var result []*Task
	for _, task := range tasks {
		match := true

		if areDone == "true" && !task.Done {
			match = false
		} else if areDone == "false" && task.Done {
			match = false
		}

		if when != "" && task.When != when {
			match = false
		}

		if match {
			result = append(result, task)
		}
	}

	json.NewEncoder(w).Encode(TasksListResponse{Data: result})
}

func getTask(w http.ResponseWriter, r *http.Request, taskID int) {
	task, ok := tasks[taskID]
	if !ok {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Task not found!"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(TaskResponse{Data: task})
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	var taskReq struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		When        string `json:"when"`
	}

	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Request must be json!"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validation
	var messages []string
	if taskReq.ID == 0 {
		messages = append(messages, "We need an id to know which entity to update!")
	}
	if taskReq.Name == "" {
		messages = append(messages, "Task name is not valid!")
	}
	if taskReq.Description == "" {
		messages = append(messages, "Task description is not valid!")
	}
	if taskReq.When == "" {
		messages = append(messages, "Task must have a date!")
	} else if _, err := time.Parse("2006-01-02", taskReq.When); err != nil {
		messages = append(messages, "`when` is not a valid date!")
	}

	if len(messages) > 0 {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Task is not valid!", Messages: messages})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, ok := tasks[taskReq.ID]
	if !ok {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Task not found!"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	task.Name = taskReq.Name
	task.Description = taskReq.Description
	task.When = taskReq.When
	task.UpdatedAt = time.Now()

	json.NewEncoder(w).Encode(TaskResponse{Data: task})
}

func completeTask(w http.ResponseWriter, r *http.Request, taskID int) {
	task, ok := tasks[taskID]
	if !ok {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Task not found!"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	task.Done = true
	task.UpdatedAt = time.Now()

	json.NewEncoder(w).Encode(TaskResponse{Data: task})
}

func deleteTask(w http.ResponseWriter, r *http.Request, taskID int) {
	task, ok := tasks[taskID]
	if !ok {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Task not found!"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(tasks, taskID)

	json.NewEncoder(w).Encode(TaskResponse{Data: task})
}