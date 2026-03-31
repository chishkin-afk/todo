package group

import (
	"testing"
	"time"

	"github.com/chishkin-afk/todo/internal/modules/task/domain/task"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGroup_Success(t *testing.T) {
	ownerID := uuid.New()
	groupID := uuid.New()
	tsk, err := task.New(ownerID, groupID, "title", "desc", task.PRIORITY_LOW)
	require.NoError(t, err)

	tasks := []*task.Task{tsk}
	title := "My Group Title"

	now := time.Now().UTC()
	group, err := New(ownerID, title, tasks)

	require.NoError(t, err)
	assert.Equal(t, ownerID, group.OwnerID())
	assert.Equal(t, title, group.Title())
	assert.Equal(t, len(tasks), len(group.Tasks()))
	assert.WithinDuration(t, now, group.CreatedAt(), 100*time.Millisecond)
	assert.WithinDuration(t, now, group.UpdatedAt(), 100*time.Millisecond)
}

func TestNewGroup_Invalid(t *testing.T) {
	ownerID := uuid.New()
	groupID := uuid.New()
	tsk, err := task.New(ownerID, groupID, "title", "desc", task.PRIORITY_LOW)
	require.NoError(t, err)

	type input struct {
		ownerID uuid.UUID
		title   string
		tasks   []*task.Task
	}

	testCases := []struct {
		name     string
		input    input
		expected error
	}{
		{
			name: "empty_title",
			input: input{
				ownerID: ownerID,
				title:   "",
				tasks:   []*task.Task{tsk},
			},
			expected: errs.ErrInvalidTitle,
		},
		{
			name: "short_title",
			input: input{
				ownerID: ownerID,
				title:   "ab",
				tasks:   []*task.Task{tsk},
			},
			expected: errs.ErrInvalidTitle,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tc.input.ownerID, tc.input.title, tc.input.tasks)

			assert.EqualError(t, err, tc.expected.Error())
		})
	}
}

func TestGroup_ChangeTitle(t *testing.T) {
	ownerID := uuid.New()
	groupID := uuid.New()
	tsk, err := task.New(ownerID, groupID, "title", "desc", task.PRIORITY_LOW)
	require.NoError(t, err)

	group, err := New(ownerID, "Old Title", []*task.Task{tsk})
	require.NoError(t, err)

	newTitle := "Updated Group Title"
	err = group.ChangeTitle(newTitle)

	require.NoError(t, err)
	assert.Equal(t, newTitle, group.Title())
	assert.NotEqual(t, group.CreatedAt(), group.UpdatedAt())
}

func TestGroup_ChangeTitle_Invalid(t *testing.T) {
	ownerID := uuid.New()
	groupID := uuid.New()
	tsk, err := task.New(ownerID, groupID, "title", "desc", task.PRIORITY_LOW)
	require.NoError(t, err)

	group, err := New(ownerID, "Valid Title", []*task.Task{tsk})
	require.NoError(t, err)

	testCases := []struct {
		name  string
		title string
		err   error
	}{
		{
			name:  "empty_title",
			title: "",
			err:   errs.ErrInvalidTitle,
		},
		{
			name:  "short_title",
			title: "ab",
			err:   errs.ErrInvalidTitle,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := group.ChangeTitle(tc.title)
			assert.EqualError(t, err, tc.err.Error())
		})
	}
}
