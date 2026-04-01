package grouppg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chishkin-afk/todo/internal/modules/task/domain/group"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type groupPersistenceRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *groupPersistenceRepository {
	return &groupPersistenceRepository{
		db: db,
	}
}

func (gpr *groupPersistenceRepository) Save(ctx context.Context, group *group.Group) (*group.Group, error) {
	_, err := gpr.db.ExecContext(ctx, `insert into groups (
		id,
		owner_id,
		title,
		created_at,
		updated_at) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)`, group.ID(), group.OwnerID(), group.Title(), group.CreatedAt(), group.UpdatedAt())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := gpr.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to save group into db: %w", err)
	}

	return group, nil
}

func (gpr *groupPersistenceRepository) GetListByUserID(ctx context.Context, userID uuid.UUID) ([]*group.Group, error) {
	rows, err := gpr.db.QueryContext(ctx, `select id, owner_id, title, created_at, updated_at from groups where owner_id = $1`, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := gpr.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get list of groups: %w", err)
	}

	return ToDomain(rows)
}

func (gpr *groupPersistenceRepository) GetByID(ctx context.Context, groupID uuid.UUID) (*group.Group, error) {
	rows, err := gpr.db.QueryContext(ctx, `
		select 
			g.id, g.owner_id, g.title, g.created_at, g.updated_at,
			t.id, t.owner_id, t.group_id, t.title, t.task_desc, t.priority_id, t.is_done, t.created_at, t.updated_at
		from groups g
		left join tasks t on g.id = t.group_id
		where g.id = $1
	`, groupID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := gpr.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get group by id: %w", err)
	}
	defer rows.Close()

	return ToDomainWithTasks(rows)
}

func (gpr *groupPersistenceRepository) Update(ctx context.Context, group *group.Group) (*group.Group, error) {
	_, err := gpr.db.ExecContext(ctx, `update groups set
		title = $1,
		updated_at = $2
		where id = $3
	`, group.Title(), group.UpdatedAt(), group.ID())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if err, ok := gpr.knownError(err); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	return group, nil
}

func (gpr *groupPersistenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := gpr.db.ExecContext(ctx, `delete from groups where id = $1`, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if err, ok := gpr.knownError(err); ok {
			return err
		}

		return fmt.Errorf("failed to delete group: %w", err)
	}

	if c, _ := result.RowsAffected(); c == 0 {
		return errs.ErrGroupNotFound
	}

	return nil
}

func (gpr *groupPersistenceRepository) knownError(err error) (error, bool) {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23503":
			return errs.ErrUserNotFound, true
		}
	}

	return err, false
}
