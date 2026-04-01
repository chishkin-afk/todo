package task

import (
	"strings"
	"time"

	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
)

type Task struct {
	id        uuid.UUID
	ownerID   uuid.UUID
	groupID   uuid.UUID
	title     string
	desc      string
	priority  priority
	isDone    bool
	createdAt time.Time
	updatedAt time.Time
}

func New(
	ownerID uuid.UUID,
	groupID uuid.UUID,
	title string,
	desc string,
	priority priority,
) (*Task, error) {
	title = strings.TrimSpace(title)
	if len([]rune(title)) < 3 || len([]rune(title)) > 64 {
		return nil, errs.ErrInvalidTitle
	}

	desc = strings.TrimSpace(desc)
	if len([]rune(desc)) < 3 || len([]rune(desc)) > 512 {
		return nil, errs.ErrInvalidTaskDesc
	}

	if !priority.IsValid() {
		return nil, errs.ErrInvalidTaskPriority
	}

	return &Task{
		id:        uuid.New(),
		ownerID:   ownerID,
		groupID:   groupID,
		title:     title,
		desc:      desc,
		priority:  priority,
		isDone:    false,
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}, nil
}

func From(
	id uuid.UUID,
	ownerID uuid.UUID,
	groupID uuid.UUID,
	title string,
	desc string,
	priority priority,
	isDone bool,
	createdAt time.Time,
	updatedAt time.Time,
) (*Task, error) {
	title = strings.TrimSpace(title)
	if len([]rune(title)) < 3 || len([]rune(title)) > 64 {
		return nil, errs.ErrInvalidTitle
	}

	desc = strings.TrimSpace(desc)
	if len([]rune(desc)) < 3 || len([]rune(desc)) > 512 {
		return nil, errs.ErrInvalidTaskDesc
	}

	if !priority.IsValid() {
		return nil, errs.ErrInvalidTaskPriority
	}

	return &Task{
		id:        id,
		ownerID:   ownerID,
		groupID:   groupID,
		title:     title,
		desc:      desc,
		priority:  priority,
		isDone:    isDone,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (t *Task) Done() error {
	if t.isDone {
		return errs.ErrTaskAlreadyDone
	}

	t.isDone = true
	t.updatedAt = time.Now().UTC()
	return nil
}

func (t *Task) NotDone() error {
	if !t.isDone {
		return errs.ErrTaskNotDone
	}

	t.isDone = false
	t.updatedAt = time.Now().UTC()
	return nil
}

func (t *Task) ChangePriority(priority priority) error {
	if !priority.IsValid() {
		return errs.ErrInvalidTaskPriority
	}

	t.priority = priority
	t.updatedAt = time.Now().UTC()
	return nil
}

func (t *Task) ChangeTitle(title string) error {
	title = strings.TrimSpace(title)
	if len([]rune(title)) < 3 || len([]rune(title)) > 64 {
		return errs.ErrInvalidTitle
	}

	t.title = title
	t.updatedAt = time.Now().UTC()
	return nil
}

func (t *Task) ChangeDesc(desc string) error {
	desc = strings.TrimSpace(desc)
	if len([]rune(desc)) < 3 || len([]rune(desc)) > 512 {
		return errs.ErrInvalidTaskDesc
	}

	t.desc = desc
	t.updatedAt = time.Now().UTC()
	return nil
}

func (t *Task) ID() uuid.UUID {
	return t.id
}

func (t *Task) GroupID() uuid.UUID {
	return t.groupID
}

func (t *Task) OwnerID() uuid.UUID {
	return t.ownerID
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Desc() string {
	return t.desc
}

func (t *Task) Priority() priority {
	return t.priority
}

func (t *Task) IsDone() bool {
	return t.isDone
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) UpdatedAt() time.Time {
	return t.updatedAt
}
