package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	mongoAdapter "github.com/lk153/quizgame-ai-serving/internal/adapters/storage/mongo"
	taskResultDomain "github.com/lk153/quizgame-ai-serving/internal/core/domains/taskResult"
	"github.com/lk153/quizgame-ai-serving/internal/core/ports"
)

const (
	collection = "task_result"
)

var _ ports.ITaskResultRepository = &TaskResultRepository{}

/**
 * TaskResultRepository implements port.TaskResultRepository interface
 * and provides an access to the postgres database
 */
type TaskResultRepository struct {
	db   *mongoAdapter.DB
	coll *mongo.Collection
}

// NewTaskResultRepository creates a task result repository instance
func NewTaskResultRepository(db *mongoAdapter.DB) *TaskResultRepository {
	coll := db.DB.Collection(collection)
	return &TaskResultRepository{
		db,
		coll,
	}
}

// Create creates a new task result in the database
func (t *TaskResultRepository) Create(
	ctx context.Context, taskResult *taskResultDomain.TaskResultEntity,
) (*taskResultDomain.TaskResultEntity, error) {
	result, err := t.coll.InsertOne(ctx, taskResult)
	if err != nil {
		return nil, err
	}

	log.Printf("Result task _id %v has been inserted", result.InsertedID)
	return taskResult, nil
}

// GetByID gets a task result by ID from the database
func (t *TaskResultRepository) GetByID(
	ctx context.Context, id string,
) (*taskResultDomain.TaskResultEntity, error) {
	var (
		taskResult taskResultDomain.TaskResultEntity
		err        error
	)
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: 1}})
	filter := bson.D{{Key: "id", Value: id}}
	if err = t.coll.FindOne(ctx, filter, opts).Decode(&taskResult); err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		log.Println("GetByID:ERROR:", err)
	}

	return &taskResult, nil
}

// List lists all task results from the database
func (t *TaskResultRepository) List(
	ctx context.Context, skip, limit uint64,
) ([]taskResultDomain.TaskResultEntity, error) {
	var tasks []taskResultDomain.TaskResultEntity
	filter := bson.D{}
	opts := options.Find().
		SetSort(bson.D{{Key: "id", Value: 1}}).
		SetLimit(int64(limit)).SetSkip(int64(skip))
	cursor, err := t.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if cursor == nil {
		return nil, nil
	}

	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Update updates a task result by ID in the database
func (t *TaskResultRepository) Update(
	ctx context.Context, task *taskResultDomain.TaskResultEntity,
) (*taskResultDomain.TaskResultEntity, error) {
	var taskResult taskResultDomain.TaskResultEntity
	opts := options.FindOneAndUpdate()
	filter := bson.D{{Key: "id", Value: task.ID}}
	err := t.coll.FindOneAndUpdate(ctx, filter, opts).Decode(&taskResult)
	if err != nil {
		return task, err
	}

	return &taskResult, nil
}

// Delete deletes a task result by ID from the database
func (t *TaskResultRepository) Delete(ctx context.Context, id string) (err error) {
	var taskResult taskResultDomain.TaskResultEntity
	opts := options.FindOneAndDelete()
	filter := bson.D{{Key: "id", Value: id}}
	err = t.coll.FindOneAndDelete(ctx, filter, opts).Decode(&taskResult)
	return
}
