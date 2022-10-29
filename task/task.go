package task

import (
	"time"

	"github.com/google/uuid"
)

type Category string
type Task struct {
	ID            uuid.UUID      `json:"id"`
	CreationTime  time.Time      `json:"creationTime"`
	ShortTitle    string         `json:"shortTitle"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Tags          []string       `json:"tags"`
	Category      Category       `json:"category"`
	Priority      uint8          `json:"priority"`
	InvolvesOther bool           `json:"involvesOther"`
	TimeEstimate  *time.Duration `json:"timeEstimate"`
	DueDate       *time.Time     `json:"dueDate"`
	HardDeadline  bool           `json:"hardDeadline"`
	Reminder      *time.Duration `json:"reminder"`
	Repeats       *time.Duration `json:"repeats"`
}

func NewEmptyTask(taskId uuid.UUID) *Task {
	return &Task{
		ID:           taskId,
		CreationTime: time.Now(),
	}
}
