package groupredis

import (
	"encoding/json"

	"github.com/chishkin-afk/todo/internal/modules/task/domain/group"
	"github.com/chishkin-afk/todo/internal/modules/task/domain/task"
)

func ToBytes(domain *group.Group) ([]byte, error) {
	model := GroupModel{
		ID:        domain.ID(),
		OwnerID:   domain.OwnerID(),
		Title:     domain.Title(),
		Tasks:     make([]TaskModel, len(domain.Tasks())),
		CreatedAt: domain.CreatedAt(),
		UpdatedAt: domain.UpdatedAt(),
	}

	for idx, task := range domain.Tasks() {
		model.Tasks[idx] = TaskModel{
			ID:         task.ID(),
			OwnerID:    task.OwnerID(),
			GroupID:    task.GroupID(),
			Title:      task.Title(),
			Desc:       task.Desc(),
			PriorityID: task.Priority().Int(),
			CreatedAt:  task.CreatedAt(),
			UpdatedAt:  task.UpdatedAt(),
		}
	}

	return json.Marshal(&model)
}

func ToDomain(bytes []byte) (*group.Group, error) {
	var model GroupModel
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	tasks := make([]*task.Task, len(model.Tasks))
	for idx, modelTask := range model.Tasks {
		priority, err := task.NewPriority(modelTask.PriorityID)
		if err != nil {
			return nil, err
		}

		task, err := task.From(
			modelTask.ID,
			modelTask.OwnerID,
			modelTask.GroupID,
			modelTask.Title,
			modelTask.Desc,
			priority,
			modelTask.CreatedAt,
			modelTask.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks[idx] = task
	}

	return group.From(
		model.ID,
		model.OwnerID,
		model.Title,
		tasks,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
