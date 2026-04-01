package services

import (
	"context"
	"log/slog"

	"github.com/chishkin-afk/todo/internal/application/dtos"
	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/chishkin-afk/todo/internal/modules/task/domain/group"
	"github.com/chishkin-afk/todo/internal/modules/task/domain/task"
	"github.com/google/uuid"
)

type taskService struct {
	cfg                  *config.Config
	log                  *slog.Logger
	taskPersistenceRepo  task.TaskPersistenceRepository
	groupPersistenceRepo group.GroupPersistenceRepository
	groupCacheRepo       group.GroupCacheRepository
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	taskPersistenceRepo task.TaskPersistenceRepository,
	groupPersistenceRepo group.GroupPersistenceRepository,
	groupCacheRepo group.GroupCacheRepository,
) *taskService {
	return &taskService{
		cfg:                  cfg,
		log:                  log,
		taskPersistenceRepo:  taskPersistenceRepo,
		groupPersistenceRepo: groupPersistenceRepo,
		groupCacheRepo:       groupCacheRepo,
	}
}

// user's id in context
func (ts *taskService) CreateGroup(ctx context.Context, req *dtos.CreateGroupRequest) (*dtos.Group, error) {
	panic("not impl")
}

// user's id in context
func (ts *taskService) UpdateGroup(ctx context.Context, req *dtos.UpdateGroupRequest) (*dtos.Group, error) {
	panic("not impl")
}

// user's id in context
func (ts *taskService) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	panic("not impl")
}

// user's id in context
func (ts *taskService) GetListGroupsByUserID(ctx context.Context) (*dtos.Groups, error) {
	panic("not impl")
}

// user's id in context
func (ts *taskService) GetGroupByID(ctx context.Context, id uuid.UUID) (*dtos.Group, error) {
	panic("not impl")
}

// user's id in context
func (ts *taskService) CreateTask(ctx context.Context, req *dtos.CreateTaskRequest) (*dtos.Task, error) {
	panic("not impl")
}

// user's id in context
func (ts *taskService) UpdateTask(ctx context.Context, req *dtos.UpdateTaskRequest) (*dtos.Task, error) {
	panic("not impl")
}

// user's id in context
func (ts *taskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	panic("not impl")
}
