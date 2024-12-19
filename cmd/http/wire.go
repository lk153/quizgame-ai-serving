//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/lk153/quizgame-ai-serving/internal/adapters/config"
	"github.com/lk153/quizgame-ai-serving/internal/adapters/http"
	"github.com/lk153/quizgame-ai-serving/internal/adapters/storage"
	mongoAdapter "github.com/lk153/quizgame-ai-serving/internal/adapters/storage/mongo"
	"github.com/lk153/quizgame-ai-serving/internal/core/services"
)

type Handlers struct {
	TaskResultHandler http.TaskResultHandler
}

var HandlerSet = wire.NewSet(
	http.NewTaskResultHandler,
	wire.Struct(new(Handlers), "TaskResultHandler"))

var SuperSet = wire.NewSet(services.ServiceSet, HandlerSet, storage.StorageSet)

func initializeDB(ctx context.Context, config *config.DB) (*mongoAdapter.DB, error) {
	return mongoAdapter.New(ctx, config)
}

func initializeHandlers(ctx context.Context, rg *gin.RouterGroup, db *mongoAdapter.DB, rd *config.Redis) Handlers {
	panic(wire.Build(SuperSet))
}
