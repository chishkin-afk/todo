package task

import (
	"context"

	"github.com/google/uuid"
)

type TaskPersistenceRepository interface {
	Save(ctx context.Context, task *Task) (*Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	Update(ctx context.Context, task *Task) (*Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
