package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriority_Success(t *testing.T) {
	priority := PRIORITY_MIDDLE

	assert.True(t, priority.IsValid())
	assert.Equal(t, "middle", priority.String())
}

func TestPriority_Invalid(t *testing.T) {
	testCases := []struct {
		name  string
		input priority
	}{
		{
			name:  "zero_priority",
			input: 0,
		},
		{
			name:  "negative_priority",
			input: -1,
		},
		{
			name:  "over_priority",
			input: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.False(t, tc.input.IsValid())
		})
	}
}
