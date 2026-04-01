package grouppg

import (
	"database/sql"
	"errors"
	"time"

	"github.com/chishkin-afk/todo/internal/modules/task/domain/group"
	"github.com/chishkin-afk/todo/internal/modules/task/domain/task"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
)

func ToDomain(rows *sql.Rows) ([]*group.Group, error) {
	groups := []*group.Group{}

	for rows.Next() {
		var id uuid.UUID
		var ownerID uuid.UUID
		var title string
		var createdAt time.Time
		var updatedAt time.Time

		if err := rows.Scan(&id, &ownerID, &title, &createdAt, &updatedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errs.ErrGroupNotFound
			}

			return nil, err
		}

		group, err := group.From(id, ownerID, title, []*task.Task{}, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrGroupNotFound
		}

		return nil, err
	}

	return groups, nil
}

func ToDomainWithTasks(rows *sql.Rows) (*group.Group, error) {
	var currentGroup *group.Group
	tasksMap := make(map[uuid.UUID]*task.Task)

	for rows.Next() {
		var (
			gID        uuid.UUID
			gOwnerID   uuid.UUID
			gTitle     string
			gCreatedAt time.Time
			gUpdatedAt time.Time

			tID        uuid.NullUUID
			tOwnerID   uuid.NullUUID
			tGroupID   uuid.NullUUID
			tTitle     sql.NullString
			tDesc      sql.NullString
			tPriority  int
			tCreatedAt sql.NullTime
			tUpdatedAt sql.NullTime
		)

		err := rows.Scan(
			&gID, &gOwnerID, &gTitle, &gCreatedAt, &gUpdatedAt,
			&tID, &tOwnerID, &tGroupID, &tTitle, &tDesc, &tPriority, &tCreatedAt, &tUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if currentGroup == nil {
			currentGroup, err = group.From(gID, gOwnerID, gTitle, []*task.Task{}, gCreatedAt, gUpdatedAt)
			if err != nil {
				return nil, err
			}
		}

		if tID.Valid {
			if _, exists := tasksMap[tID.UUID]; !exists {
				priority, err := task.NewPriority(tPriority)
				if err != nil {
					return nil, err
				}

				t, err := task.From(
					tID.UUID,
					tOwnerID.UUID,
					tGroupID.UUID,
					tTitle.String,
					tDesc.String,
					priority,
					tCreatedAt.Time,
					tUpdatedAt.Time,
				)
				if err != nil {
					return nil, err
				}
				tasksMap[tID.UUID] = t
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if currentGroup == nil {
		return nil, sql.ErrNoRows
	}

	tasks := make([]*task.Task, 0, len(tasksMap))
	for _, t := range tasksMap {
		tasks = append(tasks, t)
	}

	finalGroup, err := group.From(
		currentGroup.ID(),
		currentGroup.OwnerID(),
		currentGroup.Title(),
		tasks,
		currentGroup.CreatedAt(),
		currentGroup.UpdatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return finalGroup, nil
}
