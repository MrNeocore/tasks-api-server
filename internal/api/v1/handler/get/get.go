package get

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MrNeocore/tasks-api-server/internal/storage"
	t "github.com/MrNeocore/tasks-api-server/task"

	"github.com/google/uuid"
)

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
