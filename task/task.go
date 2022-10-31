package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Category string
type Task struct {
	ID            uuid.UUID      `json:"id" pgtype:"TEXT PRIMARY KEY"`
	CreationTime  time.Time      `json:"creationTime" pgtype:"TIMESTAMP"`
	ShortTitle    string         `json:"shortTitle" pgtype:"VARCHAR(32)"`
	Title         string         `json:"title" pgtype:"VARCHAR(256)"`
	Description   string         `json:"description" pgtype:"TEXT"`
	Tags          pq.StringArray `json:"tags" pgtype:"TEXT[]"`
	Category      Category       `json:"category" pgtype:"VARCHAR(64)"`
	Priority      uint8          `json:"priority" pgtype:"SMALLINT"`
	InvolvesOther bool           `json:"involvesOther" pgtype:"BOOL"`
	TimeEstimate  *time.Duration `json:"timeEstimate" pgtype:"INTERVAL"`
	DueDate       *time.Time     `json:"dueDate" pgtype:"TIMESTAMP"`
	HardDeadline  bool           `json:"hardDeadline" pgtype:"BOOL"`
	Reminder      *time.Duration `json:"reminder" pgtype:"INTERVAL"`
	Repeats       *time.Duration `json:"repeats" pgtype:"INTERVAL"`
}

func NewEmptyTask(taskId uuid.UUID) *Task {
	return &Task{
		ID:           taskId,
		CreationTime: time.Now(),
	}
}
