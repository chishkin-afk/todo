package group

import (
	"context"

	"github.com/google/uuid"
)

type GroupPersistenceRepository interface {
	Save(ctx context.Context, group *Group) (*Group, error)
	GetListByUserID(ctx context.Context, userID uuid.UUID) ([]*Group, error)
	GetByID(ctx context.Context, groupID uuid.UUID) (*Group, error)
	Update(ctx context.Context, group *Group) (*Group, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
