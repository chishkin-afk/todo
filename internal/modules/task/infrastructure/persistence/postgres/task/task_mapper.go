package taskpg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/chishkin-afk/todo/internal/modules/task/domain/task"
	errs "github.com/chishkin-afk/todo/pkg/errors"
)

func ToDomain(result *sql.Row) (*task.Task, error) {
	var model TaskModel
	if err := result.Scan(
		&model.ID,
		&model.OwnerID,
		&model.GroupID,
		&model.Title,
		&model.TaskDesc,
		&model.PriorityID,
		&model.IsDone,
		&model.CreatedAt,
		&model.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTaskNotFound
		}

		return nil, fmt.Errorf("failed to scan result: %w", err)
	}

	priority, err := task.NewPriority(model.PriorityID)
	if err != nil {
		return nil, err
	}

	return task.From(
		model.ID,
		model.OwnerID,
		model.GroupID,
		model.Title,
		model.TaskDesc,
		priority,
		model.IsDone,
		model.CreatedAt,
		model.UpdatedAt,
	)
}