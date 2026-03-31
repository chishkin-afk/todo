package group

import (
	"strings"
	"time"

	"github.com/chishkin-afk/todo/internal/modules/task/domain/task"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
)

type Group struct {
	id        uuid.UUID
	ownerID   uuid.UUID
	title     string
	tasks     []*task.Task
	createdAt time.Time
	updatedAt time.Time
}

func New(ownerID uuid.UUID, title string, tasks []*task.Task) (*Group, error) {
	title = strings.TrimSpace(title)
	if len([]rune(title)) < 3 || len([]rune(title)) > 64 {
		return nil, errs.ErrInvalidTitle
	}

	return &Group{
		id:        uuid.New(),
		ownerID:   ownerID,
		title:     title,
		tasks:     tasks,
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}, nil
}

func From(
	id uuid.UUID,
	ownerID uuid.UUID,
	title string,
	tasks []*task.Task,
	createdAt time.Time,
	updatedAt time.Time,
) (*Group, error) {
	title = strings.TrimSpace(title)
	if len([]rune(title)) < 3 || len([]rune(title)) > 64 {
		return nil, errs.ErrInvalidTitle
	}

	return &Group{
		id:        id,
		ownerID:   ownerID,
		title:     title,
		tasks:     tasks,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (g *Group) ChangeTitle(title string) error {
	title = strings.TrimSpace(title)
	if len([]rune(title)) < 3 || len([]rune(title)) > 64 {
		return errs.ErrInvalidTitle
	}

	g.title = title
	g.updatedAt = time.Now().UTC()
	return nil
}

func (g *Group) ID() uuid.UUID {
	return g.id
}

func (g *Group) OwnerID() uuid.UUID {
	return g.ownerID
}

func (g *Group) Title() string {
	return g.title
}

func (g *Group) Tasks() []*task.Task {
	return g.tasks
}

func (g *Group) CreatedAt() time.Time {
	return g.createdAt
}

func (g *Group) UpdatedAt() time.Time {
	return g.updatedAt
}
