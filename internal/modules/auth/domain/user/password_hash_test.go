package user

import (
	"testing"

	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPasswordHash_Success(t *testing.T) {
	password := "qwerty123"

	passwordHash, err := NewPasswordHash(password)

	require.NoError(t, err)
	assert.NotEqual(t, password, passwordHash.String())
}

func TestNewPasswordHash_Invalid(t *testing.T) {
	empty := [64]rune{}

	testCases := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "empty_password",
			input:    "",
			expected: errs.ErrInvalidPassword,
		},
		{
			name:     "too_little_password",
			input:    "",
			expected: errs.ErrInvalidPassword,
		},
		{
			name:     "too_large_password",
			input:    string(empty[:]),
			expected: errs.ErrInvalidPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewPasswordHash(tc.input)

			assert.EqualError(t, err, tc.expected.Error())
		})
	}
}
