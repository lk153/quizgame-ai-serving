package services

import (
	"github.com/google/wire"

	"github.com/lk153/quizgame-ai-serving/internal/core/ports"
	taskResultSvc "github.com/lk153/quizgame-ai-serving/internal/core/services/taskResult"
)

var ServiceSet = wire.NewSet(
	taskResultSvc.NewTaskResultService,
	wire.Bind(new(ports.ITaskResultService), new(*taskResultSvc.TaskResultService)),
)
