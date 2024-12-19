package ports

import (
	"context"

	taskResultEntities "github.com/lk153/quizgame-ai-serving/internal/core/domains/taskResult"
)

//go:generate mockgen -source=taskResult.go -destination=mocks/taskResult.go -package=mocks

// ITaskResultRepository is an interface for interacting with related task result data as CRUD
type ITaskResultRepository interface {
	// Create inserts a task result into the database
	Create(ctx context.Context, taskResult *taskResultEntities.TaskResultEntity) (*taskResultEntities.TaskResultEntity, error)

	// GetByID selects a task result by id
	GetByID(ctx context.Context, id string) (*taskResultEntities.TaskResultEntity, error)

	// List selects a list of task results with pagination
	List(ctx context.Context, skip, limit uint64) ([]taskResultEntities.TaskResultEntity, error)

	// Update updates a task result
	Update(ctx context.Context, taskResult *taskResultEntities.TaskResultEntity) (*taskResultEntities.TaskResultEntity, error)

	// Delete deletes a task result
	Delete(ctx context.Context, id string) error
}

// ITaskResultService is an interface for interacting with related task result business logic
type ITaskResultService interface {
	// SubmitTask submit task result
	SubmitTask(ctx context.Context, taskResult *taskResultEntities.TaskResultEntity) (*taskResultEntities.TaskResultEntity, error)

	// GetTaskResult returns a task result by id
	GetTaskResult(ctx context.Context, id string) (*taskResultEntities.TaskResultEntity, error)

	// ListTaskResults returns a list of task results with pagination
	ListTaskResults(ctx context.Context, skip, limit uint64) ([]taskResultEntities.TaskResultEntity, error)

	// UpdateTaskResult updates a task result
	UpdateTaskResult(ctx context.Context, user *taskResultEntities.TaskResultEntity) (*taskResultEntities.TaskResultEntity, error)

	// DeleteTaskResult deletes a task result
	DeleteTaskResult(ctx context.Context, id string) error
}
