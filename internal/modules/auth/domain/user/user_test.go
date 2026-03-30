package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser_Success(t *testing.T) {
	email := Email("mail@example.com")
	password := "qwerty123"
	username := "chishkin"

	now := time.Now().UTC()
	user, err := New(
		email,
		password,
		username,
	)

	require.NoError(t, err)
	assert.NotEmpty(t, user.id)
	assert.Equal(t, email, user.Email())
	assert.NotEqual(t, password, user.PasswordHash())
	assert.WithinDuration(t, now, user.CreatedAt(), 100*time.Millisecond)
	assert.WithinDuration(t, now, user.UpdatedAt(), 100*time.Millisecond)
}

func TestNewUser_Invalid(t *testing.T) {
	type input struct {
		email    Email
		password string
		username string
	}

	testCases := []struct {
		name     string
		input    input
		expected error
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(
				tc.input.email,
				tc.input.password,
				tc.input.username,
			)

			assert.EqualError(t, tc.expected, err.Error())
		})
	}
}
