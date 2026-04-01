package taskservices

import (
	"context"
	"errors"
	"log/slog"
	"sort"
	"time"

	"github.com/chishkin-afk/todo/internal/application/dtos"
	"github.com/chishkin-afk/todo/internal/common/config"
	"github.com/chishkin-afk/todo/internal/modules/task/domain/group"
	"github.com/chishkin-afk/todo/internal/modules/task/domain/task"
	"github.com/chishkin-afk/todo/pkg/consts"
	errs "github.com/chishkin-afk/todo/pkg/errors"
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
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := ts.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	group, err := group.New(userID, req.Title, nil)
	if err != nil {
		return nil, err
	}

	savedGroup, err := ts.groupPersistenceRepo.Save(ctx, group)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, err
		}

		ts.log.Error("failed to save group",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	return &dtos.Group{
		ID:        savedGroup.ID().String(),
		OwnerID:   savedGroup.OwnerID().String(),
		Title:     savedGroup.Title(),
		CreatedAt: savedGroup.CreatedAt().UnixMilli(),
		UpdatedAt: savedGroup.UpdatedAt().UnixMilli(),
	}, nil
}

// user's id in context
func (ts *taskService) UpdateGroup(ctx context.Context, req *dtos.UpdateGroupRequest) (*dtos.Group, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := ts.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	groupID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	group, err := ts.groupPersistenceRepo.GetByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, errs.ErrGroupNotFound) {
			return nil, err
		}

		ts.log.Error("failed to get group by id",
			slog.String("error", err.Error()),
			slog.String("group_id", groupID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	if group.OwnerID() != userID {
		return nil, errs.ErrNotEnoughRights
	}

	if err := group.ChangeTitle(req.Title); err != nil {
		return nil, err
	}

	updatedGroup, err := ts.groupPersistenceRepo.Update(ctx, group)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		ts.log.Error("failed to update group by id",
			slog.String("error", err.Error()),
			slog.String("group_id", groupID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	go ts.clearGroupCache(context.Background(), updatedGroup.ID())
	return &dtos.Group{
		ID:        updatedGroup.ID().String(),
		OwnerID:   updatedGroup.OwnerID().String(),
		Title:     updatedGroup.Title(),
		CreatedAt: updatedGroup.CreatedAt().UnixMilli(),
		UpdatedAt: updatedGroup.UpdatedAt().UnixMilli(),
	}, nil
}

// user's id in context
func (ts *taskService) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	userID, err := ts.getUserID(ctx)
	if err != nil {
		return err
	}

	group, err := ts.groupPersistenceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if errors.Is(err, errs.ErrGroupNotFound) {
			return err
		}

		ts.log.Error("failed to get group by id",
			slog.String("error", err.Error()),
			slog.String("group_id", id.String()),
		)
		return errs.ErrInternalServer
	}

	if group.OwnerID() != userID {
		return errs.ErrNotEnoughRights
	}

	if err := ts.groupPersistenceRepo.Delete(ctx, group.ID()); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if errors.Is(err, errs.ErrGroupNotFound) {
			return err
		}

		ts.log.Error("failed to delete group",
			slog.String("error", err.Error()),
			slog.String("group_id", id.String()),
		)
		return errs.ErrInternalServer
	}

	go ts.clearGroupCache(context.Background(), group.ID())
	return nil
}

// user's id in context
func (ts *taskService) GetListGroupsByUserID(ctx context.Context) (*dtos.Groups, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := ts.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := ts.groupPersistenceRepo.GetListByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		ts.log.Error("failed to get list of groups",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	response := dtos.Groups{
		Groups: make([]dtos.Group, len(groups)),
	}
	for idx, group := range groups {
		response.Groups[idx] = dtos.Group{
			ID:        group.ID().String(),
			OwnerID:   group.OwnerID().String(),
			Title:     group.Title(),
			CreatedAt: group.CreatedAt().UnixMilli(),
			UpdatedAt: group.UpdatedAt().UnixMilli(),
		}
	}

	return &response, nil
}

// user's id in context
func (ts *taskService) GetGroupByID(ctx context.Context, id uuid.UUID) (*dtos.Group, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := ts.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	if group, err := ts.groupCacheRepo.Get(ctx, id); err == nil {
		ts.log.Debug("group has been taken from cache")

		if group.OwnerID() != userID {
			return nil, errs.ErrGroupNotFound
		}

		return ts.returnGroupWithTasks(group), nil
	} else if !errors.Is(err, errs.ErrGroupNotFound) {
		ts.log.Error("failed to get group cache",
			slog.String("error", err.Error()),
			slog.String("group_id", id.String()),
		)
	}

	group, err := ts.groupPersistenceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, errs.ErrGroupNotFound) {
			return nil, err
		}

		ts.log.Error("failed to get group with tasks by id",
			slog.String("error", err.Error()),
			slog.String("group_id", id.String()),
		)
		return nil, errs.ErrInternalServer
	}

	if group.OwnerID() != userID {
		return nil, errs.ErrGroupNotFound
	}

	go ts.saveGroupCache(context.Background(), group)
	return ts.returnGroupWithTasks(group), nil
}

func (ts *taskService) returnGroupWithTasks(group *group.Group) *dtos.Group {
	response := dtos.Group{
		ID:        group.ID().String(),
		OwnerID:   group.OwnerID().String(),
		Title:     group.Title(),
		Tasks:     make([]dtos.Task, len(group.Tasks())),
		CreatedAt: group.CreatedAt().UnixMilli(),
		UpdatedAt: group.UpdatedAt().UnixMilli(),
	}
	for idx, task := range group.Tasks() {
		response.Tasks[idx] = dtos.Task{
			ID:         task.ID().String(),
			OwnerID:    task.OwnerID().String(),
			GroupID:    task.GroupID().String(),
			Title:      task.Title(),
			Desc:       task.Desc(),
			Priority:   task.Priority().String(),
			PriorityID: int64(task.Priority().Int()),
			IsDone:     task.IsDone(),
			CreatedAt:  task.CreatedAt().UnixMilli(),
			UpdatedAt:  task.UpdatedAt().UnixMilli(),
		}
	}

	sort.Slice(response.Tasks, func(i, j int) bool {
		return response.Tasks[i].PriorityID > response.Tasks[j].PriorityID
	})

	return &response
}

// user's id in context
func (ts *taskService) CreateTask(ctx context.Context, req *dtos.CreateTaskRequest) (*dtos.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := ts.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	priority, err := task.NewPriority(req.PriorityID)
	if err != nil {
		return nil, err
	}

	groupID, err := uuid.Parse(req.GroupID)
	if err != nil {
		return nil, errs.ErrGroupNotFound
	}

	task, err := task.New(userID, groupID, req.Title, req.Desc, priority)
	if err != nil {
		return nil, err
	}

	savedTask, err := ts.taskPersistenceRepo.Save(ctx, task)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		ts.log.Error("failed to save task",
			slog.String("error", err.Error()),
		)
		return nil, errs.ErrInternalServer
	}

	go ts.clearGroupCache(context.Background(), groupID)
	return &dtos.Task{
		ID:         savedTask.ID().String(),
		OwnerID:    savedTask.OwnerID().String(),
		GroupID:    savedTask.GroupID().String(),
		Title:      savedTask.Title(),
		Desc:       savedTask.Desc(),
		Priority:   savedTask.Priority().String(),
		PriorityID: int64(savedTask.Priority().Int()),
		IsDone:     savedTask.IsDone(),
		CreatedAt:  savedTask.CreatedAt().UnixMilli(),
		UpdatedAt:  savedTask.UpdatedAt().UnixMilli(),
	}, nil
}

// user's id in context
func (ts *taskService) UpdateTask(ctx context.Context, req *dtos.UpdateTaskRequest) (*dtos.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userID, err := ts.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	taskID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	taskToUpdate, err := ts.taskPersistenceRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		if errors.Is(err, errs.ErrTaskNotFound) {
			return nil, err
		}

		ts.log.Error("failed to get task by id",
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	if taskToUpdate.OwnerID() != userID {
		return nil, errs.ErrNotEnoughRights
	}

	if err = ts.applyUpdates(taskToUpdate, req); err != nil {
		return nil, err
	}

	updatedTask, err := ts.taskPersistenceRepo.Update(ctx, taskToUpdate)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		ts.log.Error("failed to update task",
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
		)
		return nil, errs.ErrInternalServer
	}

	go ts.clearGroupCache(ctx, updatedTask.GroupID())
	return &dtos.Task{
		ID:         updatedTask.ID().String(),
		OwnerID:    updatedTask.OwnerID().String(),
		GroupID:    updatedTask.GroupID().String(),
		Title:      updatedTask.Title(),
		Desc:       updatedTask.Desc(),
		Priority:   updatedTask.Priority().String(),
		PriorityID: int64(updatedTask.Priority().Int()),
		IsDone:     updatedTask.IsDone(),
		CreatedAt:  updatedTask.CreatedAt().UnixMilli(),
		UpdatedAt:  updatedTask.UpdatedAt().UnixMilli(),
	}, nil
}

// user's id in context
func (ts *taskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	userID, err := ts.getUserID(ctx)
	if err != nil {
		return err
	}

	task, err := ts.taskPersistenceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if errors.Is(err, errs.ErrTaskNotFound) {
			return err
		}

		ts.log.Error("failed to get task by id",
			slog.String("error", err.Error()),
			slog.String("task_id", id.String()),
		)
		return errs.ErrInternalServer
	}

	if task.OwnerID() != userID {
		return errs.ErrNotEnoughRights
	}

	if err := ts.taskPersistenceRepo.Delete(ctx, task.ID()); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if errors.Is(err, errs.ErrTaskNotFound) {
			return err
		}

		ts.log.Error("failed to get task by id",
			slog.String("error", err.Error()),
			slog.String("task_id", id.String()),
		)
		return errs.ErrInternalServer
	}

	go ts.clearGroupCache(ctx, task.GroupID())
	return nil
}

func (ts *taskService) getUserID(ctx context.Context) (uuid.UUID, error) {
	raw := ctx.Value(consts.UserID)
	if raw == nil {
		return uuid.Nil, errs.ErrInvalidToken
	}

	if id, ok := raw.(uuid.UUID); ok {
		return id, nil
	}

	return uuid.Nil, errs.ErrInvalidToken
}

func (ts *taskService) clearGroupCache(ctx context.Context, id uuid.UUID) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := ts.groupCacheRepo.Del(ctxTimeout, id); err != nil {
		ts.log.Error("failed to delete group from cache",
			slog.String("error", err.Error()),
			slog.String("group_id", id.String()),
		)
	}
}

func (ts *taskService) saveGroupCache(ctx context.Context, group *group.Group) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := ts.groupCacheRepo.Save(ctxTimeout, group); err != nil {
		ts.log.Error("failed to save group into cache",
			slog.String("error", err.Error()),
		)
	}
}

func (ts *taskService) applyUpdates(taskToUpdate *task.Task, req *dtos.UpdateTaskRequest) error {
	if req.Title != nil {
		if err := taskToUpdate.ChangeTitle(*req.Title); err != nil {
			return err
		}
	}

	if req.Desc != nil {
		if err := taskToUpdate.ChangeDesc(*req.Desc); err != nil {
			return err
		}
	}

	if req.IsDone != nil {
		var err error
		if *req.IsDone {
			err = taskToUpdate.Done()
		} else {
			err = taskToUpdate.NotDone()
		}
		if err != nil {
			return err
		}
	}

	if req.PriorityID != nil {
		priority, err := task.NewPriority(*req.PriorityID)
		if err != nil {
			return err
		}

		if err := taskToUpdate.ChangePriority(priority); err != nil {
			return err
		}
	}

	return nil
}
