package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmail_Success(t *testing.T) {
	email := "mail@example.com"
	emailVO := Email(email)

	require.True(t, emailVO.IsValid())
}

func TestEmail_Invalid(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "empty_email",
			input: "",
		},
		{
			name:  "empty_domain",
			input: "mail@",
		},
		{
			name:  "empty_name",
			input: "@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			email := Email(tc.input)

			require.False(t, email.IsValid())
		})
	}
}
