package store

import (
	"context"
	"fmt"

	"github.com/MrNeocore/tasks-api-server/internal/storage"
	t "github.com/MrNeocore/tasks-api-server/task"
	"github.com/lib/pq"
)

func StoreTask(ctx context.Context, task t.Task) error {
	fmt.Printf("Storing task %v.\n", task)

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
		ctx,
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
		return insertErr
	}

	fmt.Printf("Task %v stored !\n", task.ID)

	return nil
}
