package taskpg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chishkin-afk/todo/internal/modules/task/domain/task"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type taskPersistenceRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *taskPersistenceRepository {
	return &taskPersistenceRepository{
		db: db,
	}
}

func (tpr *taskPersistenceRepository) Save(ctx context.Context, task *task.Task) (*task.Task, error) {
	_, err := tpr.db.ExecContext(ctx, `insert into tasks (
		id,
		owner_id,
		group_id,
		title,
		task_desc,
		priority_id,
		created_at,
		updated_at
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8
	)`, task.ID(), task.OwnerID(), task.GroupID(), task.Title(), task.Desc(), task.Priority().Int(), task.CreatedAt(), task.UpdatedAt())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := tpr.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to save task into db: %w", err)
	}

	return task, nil
}

func (tpr *taskPersistenceRepository) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	result := tpr.db.QueryRowContext(ctx, `select 
		id,
		owner_id,
		group_id,
		title,
		task_desc,
		priority_id,
		created_at,
		updated_at
	from tasks where id = $1`, id)
	if err := result.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := tpr.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("faield to get task by id: %w", err)
	}

	return ToDomain(result)
}

func (tpr *taskPersistenceRepository) Update(ctx context.Context, task *task.Task) (*task.Task, error) {
	_, err := tpr.db.ExecContext(ctx, `update tasks set
		priority_id = $1,
		title = $2,
		task_desc = $3
		where id = $4`, task.Priority().Int(), task.Title(), task.Desc(), task.ID())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := tpr.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

func (tpr *taskPersistenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := tpr.db.ExecContext(ctx, `delete from tasks where id = $1`, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if err, ok := tpr.knownError(err); ok {
			return err
		}

		return fmt.Errorf("failed to delete task: %w", err)
	}

	if count, _ := result.RowsAffected(); count == 0 {
		return errs.ErrTaskNotFound
	}

	return nil
}

func (tpr *taskPersistenceRepository) knownError(err error) (error, bool) {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23503":
			return errs.ErrDepsNotFound, true
		}
	}

	return err, false
}
