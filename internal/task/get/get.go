package get

import (
	"context"
	"fmt"

	"github.com/MrNeocore/tasks-api-server/internal/storage"
	t "github.com/MrNeocore/tasks-api-server/task"

	"github.com/google/uuid"
)

func GetTask(ctx context.Context, taskId uuid.UUID) (*t.Task, error) {
	fmt.Printf("Getting task %v\n", taskId)

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

	fmt.Printf("Task (id: %v) fetched.\n", taskId)

	return task, nil
}

func GetTasks(ctx context.Context) ([]t.Task, error) {
	fmt.Println("Getting tasks")

	selectStmt := `SELECT
			id,
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
	`

	rows, queryErr := storage.DB.QueryContext(ctx, selectStmt)

	if queryErr != nil {
		return nil, queryErr
	}

	defer rows.Close()

	var tasks []t.Task

	for rows.Next() {
		var task t.Task
		scanErr := rows.Scan(
			&task.ID,
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

		tasks = append(tasks, task)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return tasks, nil
}
