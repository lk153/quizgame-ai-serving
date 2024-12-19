package taskresult

import (
	"context"
	"log"

	errDomain "github.com/lk153/quizgame-ai-serving/internal/core/domains/error"
	taskResultEntities "github.com/lk153/quizgame-ai-serving/internal/core/domains/taskResult"
	"github.com/lk153/quizgame-ai-serving/internal/core/ports"
	cacheLib "github.com/lk153/quizgame-ai-serving/lib/cache"
	errLib "github.com/lk153/quizgame-ai-serving/lib/errors"
)

var (
	_               ports.ITaskResultService = &TaskResultService{}
	cachePrefix                              = "taskResult"
	cacheListPrefix                          = "taskResults"
)

type TaskResultService struct {
	repo  ports.ITaskResultRepository
	cache ports.ICacheRepository
}

func NewTaskResultService(repo ports.ITaskResultRepository, cache ports.ICacheRepository) *TaskResultService {
	return &TaskResultService{
		repo,
		cache,
	}
}

// Register: create a new task result
func (u *TaskResultService) SubmitTask(
	ctx context.Context, task *taskResultEntities.TaskResultEntity,
) (e *taskResultEntities.TaskResultEntity, err error) {
	var (
		cacheKey       string
		taskSerialized []byte
	)

	task, err = u.repo.Create(ctx, task)
	if err != nil {
		errLib.Error.Println(err)
		if err == errDomain.ErrConflictingData {
			return
		}

		err = errDomain.ErrInternal
		return
	}

	cacheKey = cacheLib.GenerateCacheKey(cachePrefix, task.ID)
	taskSerialized, err = cacheLib.Serialize(task)
	if err != nil {
		goto ERR
	}

	if err = u.cache.Set(ctx, cacheKey, taskSerialized, 0); err != nil {
		goto ERR
	}

	if err = u.cache.DeleteByPrefix(ctx, cacheListPrefix+":*"); err != nil {
		goto ERR
	}

	return

ERR:
	errLib.Error.Println(err)
	err = errDomain.ErrInternal
	return
}

// GetTaskResult: return a task result by id
func (u *TaskResultService) GetTaskResult(
	ctx context.Context, id string,
) (e *taskResultEntities.TaskResultEntity, err error) {
	cacheKey := cacheLib.GenerateCacheKey(cachePrefix, id)
	cachedTask, err := u.cache.Get(ctx, cacheKey)
	if err == nil {
		err = cacheLib.Deserialize(cachedTask, &e)
		if err != nil {
			err = errDomain.ErrInternal
		}

		return
	}

	e, err = u.repo.GetByID(ctx, id)
	if err != nil {
		if err == errDomain.ErrDataNotFound {
			return
		}

		err = errDomain.ErrInternal
		return
	}

	taskSerialized, err := cacheLib.Serialize(e)
	if err != nil {
		err = errDomain.ErrInternal
		return
	}

	err = u.cache.Set(ctx, cacheKey, taskSerialized, 0)
	if err != nil {
		err = errDomain.ErrInternal
	}

	return
}

// ListTaskResults: return a list of task results with pagination
func (u *TaskResultService) ListTaskResults(
	ctx context.Context, skip, limit uint64,
) (tasks []taskResultEntities.TaskResultEntity, err error) {
	params := cacheLib.GenerateCacheKeyParams(skip, limit)
	cacheKey := cacheLib.GenerateCacheKey(cacheListPrefix, params)
	cachedTasks, err := u.cache.Get(ctx, cacheKey)
	if err == nil {
		err = cacheLib.Deserialize(cachedTasks, &tasks)
		if err != nil {
			log.Println("ERR:", err)
			err = errDomain.ErrInternal
		}

		return
	}

	tasks, err = u.repo.List(ctx, skip, limit)
	if err != nil {
		log.Println("ERR:", err)
		err = errDomain.ErrInternal
		return
	}

	tasksSerialized, err := cacheLib.Serialize(tasks)
	if err != nil {
		log.Println("ERR:", err)
		err = errDomain.ErrInternal
		return
	}

	err = u.cache.Set(ctx, cacheKey, tasksSerialized, 0)
	if err != nil {
		log.Println("ERR:", err)
		err = errDomain.ErrInternal
	}

	return
}

// UpdateTaskResult: update a task result
func (u *TaskResultService) UpdateTaskResult(
	ctx context.Context, task *taskResultEntities.TaskResultEntity,
) (e *taskResultEntities.TaskResultEntity, err error) {
	existingTask, err := u.repo.GetByID(ctx, task.ID)
	if err != nil {
		if err == errDomain.ErrDataNotFound {
			return nil, err
		}
		return nil, errDomain.ErrInternal
	}

	emptyData := task.Name == ""
	sameData := existingTask.Name == task.Name
	if emptyData || sameData {
		return nil, errDomain.ErrNoUpdatedData
	}

	_, err = u.repo.Update(ctx, task)
	if err != nil {
		if err == errDomain.ErrConflictingData {
			return
		}

		err = errDomain.ErrInternal
		return
	}

	cacheKey := cacheLib.GenerateCacheKey(cachePrefix, task.ID)
	if err = u.cache.Delete(ctx, cacheKey); err != nil {
		err = errDomain.ErrInternal
		return
	}

	taskSerialized, err := cacheLib.Serialize(task)
	if err != nil {
		err = errDomain.ErrInternal
		return
	}

	if err = u.cache.Set(ctx, cacheKey, taskSerialized, 0); err != nil {
		err = errDomain.ErrInternal
		return
	}

	if err = u.cache.DeleteByPrefix(ctx, cacheListPrefix+":*"); err != nil {
		err = errDomain.ErrInternal
	}

	return
}

// DeleteTaskResult: delete a task result
func (u *TaskResultService) DeleteTaskResult(ctx context.Context, id string) (err error) {
	_, err = u.repo.GetByID(ctx, id)
	if err != nil {
		if err == errDomain.ErrDataNotFound {
			return
		}

		return errDomain.ErrInternal
	}

	cacheKey := cacheLib.GenerateCacheKey(cachePrefix, id)
	if err = u.cache.Delete(ctx, cacheKey); err != nil {
		return errDomain.ErrInternal
	}

	if err = u.cache.DeleteByPrefix(ctx, cacheListPrefix+":*"); err != nil {
		return errDomain.ErrInternal
	}

	return u.repo.Delete(ctx, id)
}
