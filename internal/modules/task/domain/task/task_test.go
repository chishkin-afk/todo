package task

import (
	"testing"
	"time"

	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTask_Success(t *testing.T) {
	title := "Wash dishes"
	desc := "take the dishes and put it under water"
	priority := PRIORITY_HIGH
	ownerID := uuid.New()
	groupID := uuid.New()

	now := time.Now().UTC()
	task, err := New(ownerID, groupID, title, desc, priority)

	require.NoError(t, err)
	assert.Equal(t, title, task.Title())
	assert.Equal(t, desc, task.Desc())
	assert.Equal(t, priority, task.Priority())
	assert.False(t, task.IsDone())
	assert.Equal(t, ownerID, task.OwnerID())
	assert.Equal(t, groupID, task.GroupID())
	assert.WithinDuration(t, now, task.CreatedAt(), 100*time.Millisecond)
	assert.WithinDuration(t, now, task.UpdatedAt(), 100*time.Millisecond)
}

func TestNewTask_Invalid(t *testing.T) {
	empty := [1024]rune{}
	ownerID := uuid.New()
	groupID := uuid.New()

	type input struct {
		title    string
		desc     string
		priority priority
	}

	testCases := []struct {
		name     string
		input    input
		expected error
	}{
		{
			name: "empty_title",
			input: input{
				title:    "",
				desc:     "any description",
				priority: PRIORITY_LOW,
			},
			expected: errs.ErrInvalidTitle,
		},
		{
			name: "too_large_title",
			input: input{
				title:    string(empty[:]),
				desc:     "any description",
				priority: PRIORITY_LOW,
			},
			expected: errs.ErrInvalidTitle,
		},
		{
			name: "empty_desc",
			input: input{
				title:    "title",
				desc:     "",
				priority: PRIORITY_LOW,
			},
			expected: errs.ErrInvalidTaskDesc,
		},
		{
			name: "too_large_desc",
			input: input{
				title:    "title",
				desc:     string(empty[:]),
				priority: PRIORITY_LOW,
			},
			expected: errs.ErrInvalidTaskDesc,
		},
		{
			name: "invalid_priority",
			input: input{
				title:    "title",
				desc:     "any descriprion",
				priority: 0,
			},
			expected: errs.ErrInvalidTaskPriority,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(
				ownerID,
				groupID,
				tc.input.title,
				tc.input.desc,
				tc.input.priority,
			)

			assert.EqualError(t, err, tc.expected.Error())
		})
	}
}

func TestFromTask_Success(t *testing.T) {
	id := uuid.New()
	ownerID := uuid.New()
	groupID := uuid.New()
	title := "Existing Task"
	desc := "Description of existing task"
	priority := PRIORITY_MIDDLE
	isDone := true
	createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 14, 0, 0, 0, time.UTC)

	task, err := From(id, ownerID, groupID, title, desc, priority, isDone, createdAt, updatedAt)

	require.NoError(t, err)
	assert.Equal(t, id, task.ID())
	assert.Equal(t, ownerID, task.OwnerID())
	assert.Equal(t, groupID, task.GroupID())
	assert.Equal(t, title, task.Title())
	assert.Equal(t, desc, task.Desc())
	assert.Equal(t, priority, task.Priority())
	assert.Equal(t, isDone, task.IsDone())
	assert.Equal(t, createdAt, task.CreatedAt())
	assert.Equal(t, updatedAt, task.UpdatedAt())
}

func TestFromTask_Invalid(t *testing.T) {
	id := uuid.New()
	ownerID := uuid.New()
	groupID := uuid.New()

	testCases := []struct {
		name     string
		title    string
		desc     string
		priority priority
		expected error
	}{
		{
			name:     "invalid_title",
			title:    "ab",
			desc:     "valid description",
			priority: PRIORITY_LOW,
			expected: errs.ErrInvalidTitle,
		},
		{
			name:     "invalid_desc",
			title:    "valid title",
			desc:     "ab",
			priority: PRIORITY_LOW,
			expected: errs.ErrInvalidTaskDesc,
		},
		{
			name:     "invalid_priority",
			title:    "valid title",
			desc:     "valid description",
			priority: PRIORITY_UNKNOWN,
			expected: errs.ErrInvalidTaskPriority,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := From(id, ownerID, groupID, tc.title, tc.desc, tc.priority, false, time.Now(), time.Now())
			assert.EqualError(t, err, tc.expected.Error())
		})
	}
}

func TestTask_ChangePriority(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)

	err = task.ChangePriority(PRIORITY_HIGH)
	require.NoError(t, err)
	assert.Equal(t, PRIORITY_HIGH, task.Priority())
	assert.NotEqual(t, task.CreatedAt(), task.UpdatedAt())
}

func TestTask_ChangePriority_Invalid(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)

	err = task.ChangePriority(PRIORITY_UNKNOWN)
	assert.EqualError(t, err, errs.ErrInvalidTaskPriority.Error())
}

func TestTask_ChangeTitle(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)

	newTitle := "Updated Title"
	err = task.ChangeTitle(newTitle)
	require.NoError(t, err)
	assert.Equal(t, newTitle, task.Title())
	assert.NotEqual(t, task.CreatedAt(), task.UpdatedAt())
}

func TestTask_ChangeTitle_Invalid(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)

	err = task.ChangeTitle("ab")
	assert.EqualError(t, err, errs.ErrInvalidTitle.Error())
}

func TestTask_ChangeDesc(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)

	newDesc := "Updated Description"
	err = task.ChangeDesc(newDesc)
	require.NoError(t, err)
	assert.Equal(t, newDesc, task.Desc())
	assert.NotEqual(t, task.CreatedAt(), task.UpdatedAt())
}

func TestTask_ChangeDesc_Invalid(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)

	err = task.ChangeDesc("ab")
	assert.EqualError(t, err, errs.ErrInvalidTaskDesc.Error())
}

func TestTask_Done(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)
	require.False(t, task.IsDone())

	err = task.Done()
	require.NoError(t, err)
	assert.True(t, task.IsDone())
	assert.NotEqual(t, task.CreatedAt(), task.UpdatedAt())
}

func TestTask_Done_AlreadyDone(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)
	_ = task.Done()

	err = task.Done()
	assert.EqualError(t, err, errs.ErrTaskAlreadyDone.Error())
}

func TestTask_NotDone(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)
	_ = task.Done()
	require.True(t, task.IsDone())

	err = task.NotDone()
	require.NoError(t, err)
	assert.False(t, task.IsDone())
	assert.NotEqual(t, task.CreatedAt(), task.UpdatedAt())
}

func TestTask_NotDone_AlreadyNotDone(t *testing.T) {
	task, err := New(uuid.New(), uuid.New(), "Test Title", "Test Desc", PRIORITY_LOW)
	require.NoError(t, err)
	require.False(t, task.IsDone())

	err = task.NotDone()
	assert.EqualError(t, err, errs.ErrTaskNotDone.Error())
}
