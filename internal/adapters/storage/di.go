package storage

import (
	"context"
	"fmt"

	"github.com/google/wire"

	"github.com/lk153/quizgame-ai-serving/internal/adapters/config"
	"github.com/lk153/quizgame-ai-serving/internal/adapters/storage/mongo/repository"
	"github.com/lk153/quizgame-ai-serving/internal/adapters/storage/redis"
	"github.com/lk153/quizgame-ai-serving/internal/core/ports"
)

func ProvideRedis(ctx context.Context, config *config.Redis) *redis.Redis {
	rd, err := redis.New(ctx, config)
	if err != nil {
		panic(fmt.Sprintf("Connect redis failed: %s", err.Error()))
	}

	return rd
}

var StorageSet = wire.NewSet(
	repository.NewTaskResultRepository,
	wire.Bind(new(ports.ITaskResultRepository), new(*repository.TaskResultRepository)),

	ProvideRedis,
	wire.Bind(new(ports.ICacheRepository), new(*redis.Redis)),
)
