package groupredis

import (
	"time"

	"github.com/google/uuid"
)

type GroupModel struct {
	ID        uuid.UUID   `json:"id"`
	OwnerID   uuid.UUID   `json:"owner_id"`
	Title     string      `json:"title"`
	Tasks     []TaskModel `json:"tasks"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type TaskModel struct {
	ID         uuid.UUID `json:"id"`
	OwnerID    uuid.UUID `json:"owner_id"`
	GroupID    uuid.UUID `json:"group_id"`
	Title      string    `json:"title"`
	Desc       string    `json:"desc"`
	PriorityID int       `json:"priority_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
