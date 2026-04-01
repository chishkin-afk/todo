package taskpg

import (
	"time"

	"github.com/google/uuid"
)

type TaskModel struct {
	ID         uuid.UUID
	OwnerID    uuid.UUID
	GroupID    uuid.UUID
	Title      string
	TaskDesc   string
	PriorityID int
	IsDone     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}