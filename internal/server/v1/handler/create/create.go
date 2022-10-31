package create

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MrNeocore/tasks-api-server/internal/storage"
	t "github.com/MrNeocore/tasks-api-server/task"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

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
