package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang-rest-api-server/internal/storage"
	t "golang-rest-api-server/task"

	"github.com/google/uuid"
)

func Tasks(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		CreateTask(w, req)
	case http.MethodGet:
		pathParts := strings.Split(req.URL.Path, "/")
		if len(pathParts) != 3 {
			GetTasks(w, req)
		} else {
			_taskId := pathParts[2]
			GetTask(_taskId, w, req)
		}
	default:
		errorMessage := "Unsupported verb for route."
		fmt.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
	}
}

func CreateTask(w http.ResponseWriter, req *http.Request) {
	task, taskCreationError := newTaskFromRequest(req)

	if taskCreationError != nil {
		http.Error(w, taskCreationError.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Adding task %v.", task)

	_, insertErr := storage.DB.ExecContext(req.Context(), "INSERT INTO tasks (id, title) VALUES ($1, $2);", task.ID, task.Title)

	if insertErr != nil {
		errorMessage := "Failed to store the task !"
		fmt.Println(errorMessage, insertErr.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	successMessage := fmt.Sprintf("Task %v added\n", task.ID)

	fmt.Print(successMessage)
	fmt.Fprint(w, successMessage)
}

func newTaskFromRequest(req *http.Request) (*t.Task, error) {
	taskId := uuid.New().String()

	task := t.Task{ID: taskId}
	decodeErr := json.NewDecoder(req.Body).Decode(&task)

	if decodeErr != nil {
		return nil, decodeErr
	}


func GetTasks(w http.ResponseWriter, req *http.Request) {
	tasks, getTasksError := _getTasks(req.Context())

	if getTasksError != nil {
		fmt.Println(getTasksError.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(tasks) == 0 {
		noTasksMessage := "No task exist yet."
		fmt.Println(noTasksMessage)
		http.Error(w, noTasksMessage, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	marshalError := json.NewEncoder(w).Encode(tasks)

	if marshalError != nil {
		fmt.Println(marshalError.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func _getTasks(ctx context.Context) ([]t.Task, error) {
	rows, queryErr := storage.DB.QueryContext(ctx, "SELECT ID, title FROM tasks")

	if queryErr != nil {
		return nil, queryErr
	}

	defer rows.Close()

	var tasks []t.Task

	for rows.Next() {
		var task t.Task
		if err := rows.Scan(&task.ID, &task.Title); err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
}

	if rowsErr := rows.Err(); rowsErr != nil {
		return tasks, rowsErr
	}

	return tasks, nil
}

func GetTask(_taskId string, w http.ResponseWriter, req *http.Request) {
	taskId, taskIdParseError := uuid.Parse(_taskId)
	if taskIdParseError != nil {
		taskIdParseErrorMessage := fmt.Sprintf("Failed to parse taskId: %v\n", _taskId)
		fmt.Print(taskIdParseErrorMessage)
		http.Error(w, taskIdParseErrorMessage, http.StatusBadRequest)
		return
	}

	task, getTaskError := _getTask(req.Context(), taskId)

	if getTaskError != nil {
		fmt.Println(getTaskError.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if task == nil {
		taskNotFoundErrorMessage := fmt.Sprintf("Task ID %v not found", taskId)
		fmt.Println(taskNotFoundErrorMessage)
		http.Error(w, taskNotFoundErrorMessage, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	marshalError := json.NewEncoder(w).Encode(task)

	if marshalError != nil {
		fmt.Println(marshalError.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func _getTask(ctx context.Context, taskId uuid.UUID) (*t.Task, error) {
	rows, queryErr := storage.DB.QueryContext(ctx, "SELECT title FROM tasks WHERE id = $1", taskId)

	if queryErr != nil {
		return nil, queryErr
	}

	defer rows.Close()

	task := &t.Task{ID: taskId}

	// Task not found
	if !rows.Next() {
		return nil, nil
	}

	scanErr := rows.Scan(&task.Title)
	if scanErr != nil {
		return nil, scanErr
	}

	return task, nil
}
