package task

import errs "github.com/chishkin-afk/todo/pkg/errors"

type priority int

const (
	PRIORITY_UNKNOWN priority = iota
	PRIORITY_LOW
	PRIORITY_MIDDLE
	PRIORITY_HIGH
)

func (p priority) String() string {
	switch p {
	case PRIORITY_LOW:
		return "low"
	case PRIORITY_MIDDLE:
		return "middle"
	case PRIORITY_HIGH:
		return "high"
	}

	return "unknown"
}

func (p priority) Int() int {
	return int(p)
}

func (p priority) IsValid() bool {
	if p <= PRIORITY_UNKNOWN || p > PRIORITY_HIGH {
		return false
	}

	return true
}

func NewPriority(val int) (priority, error) {
	p := priority(val)
	if !p.IsValid() {
		return PRIORITY_UNKNOWN, errs.ErrInvalidTaskPriority
	}
	return p, nil
}
