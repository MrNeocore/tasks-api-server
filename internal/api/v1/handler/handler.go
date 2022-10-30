package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/MrNeocore/tasks-api-server/internal/storage"
	t "github.com/MrNeocore/tasks-api-server/task"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func Tasks(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		CreateTask(w, req)
	case http.MethodGet:
		_taskId, hasValidTaskId, taskIdParsingError := extractValidTaskIdFromPath(req.URL)
		if hasValidTaskId {
			if taskIdParsingError != nil {
				handleTaskIdParsingError(w, taskIdParsingError)
				return
			}
			GetTask(*_taskId, w, req)
		} else {
			GetTasks(w, req)
		}
	default:
		errorMessage := "Unsupported verb for route."
		fmt.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
	}
}

func extractValidTaskIdFromPath(url *url.URL) (*uuid.UUID, bool, error) {
	pathParts := strings.Split(url.Path, "/")
	if len(pathParts) != 3 {
		return nil, false, nil
	} else {
		_taskId := pathParts[2]
		taskId, taskIdParseError := uuid.Parse(_taskId)
		if taskIdParseError != nil {
			return &taskId, true, taskIdParseError
		} else {
			return &taskId, true, nil
		}
	}
}

func handleTaskIdParsingError(w http.ResponseWriter, parsingError error) {
	taskIdParseErrorMessage := fmt.Sprintf("Failed to parse taskId: %v", parsingError)
	fmt.Println(taskIdParseErrorMessage)
	http.Error(w, taskIdParseErrorMessage, http.StatusBadRequest)
}

func CreateTask(w http.ResponseWriter, req *http.Request) {
	task, taskCreationError := newTaskFromRequest(req)

	if taskCreationError != nil {
		http.Error(w, taskCreationError.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Adding task %v.\n", task)

	sqlStatement := `INSERT INTO tasks (id, 
									   creationTime, 
									   shortTitle, 
									   title, 
									   description, 
									   tags, 
									   category, 
									   priority, 
									   involvesOther, 
									   timeEstimate, 
									   dueDate, 
									   hardDeadline, 
									   reminder, 
									   repeats)

					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
	`
	_, insertErr := storage.DB.ExecContext(
		req.Context(),
		sqlStatement,
		task.ID,
		task.CreationTime,
		task.ShortTitle,
		task.Title,
		task.Description,
		pq.StringArray(task.Tags),
		task.Category,
		task.Priority,
		task.InvolvesOther,
		task.TimeEstimate,
		task.DueDate,
		task.HardDeadline,
		task.Reminder,
		task.Repeats,
	)

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
	taskId := uuid.New()

	task := t.NewEmptyTask(taskId)
	decodeErr := json.NewDecoder(req.Body).Decode(task)

	if decodeErr != nil {
		return nil, decodeErr
	}

	return task, nil
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

func GetTask(taskId uuid.UUID, w http.ResponseWriter, req *http.Request) {
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
	fmt.Printf("Fetching task %v\n", taskId)

	selectStmt := `SELECT
			creationTime,
			shortTitle,
			title,
			description,
			tags,
			category,
			priority,
			involvesOther,
			timeEstimate,
			dueDate,
			hardDeadline,
			reminder,
			repeats
		FROM tasks 
		WHERE id = $1;
	`

	rows, queryErr := storage.DB.QueryContext(ctx, selectStmt, taskId)

	if queryErr != nil {
		return nil, queryErr
	}

	defer rows.Close()

	task := &t.Task{ID: taskId}

	// Task not found
	if !rows.Next() {
		return nil, nil
	}

	scanErr := rows.Scan(
		&task.CreationTime,
		&task.ShortTitle,
		&task.Title,
		&task.Description,
		&task.Tags,
		&task.Category,
		&task.Priority,
		&task.InvolvesOther,
		&task.TimeEstimate,
		&task.DueDate,
		&task.HardDeadline,
		&task.Reminder,
		&task.Repeats,
	)

	if scanErr != nil {
		return nil, scanErr
	}

	fmt.Printf("Task %v fetched.\n", taskId)

	return task, nil
}
