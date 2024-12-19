package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	taskResultDomain "github.com/lk153/quizgame-ai-serving/internal/core/domains/taskResult"
	"github.com/lk153/quizgame-ai-serving/internal/core/ports"
)

// TaskResultHandler represents the HTTP handler for related task result requests
type TaskResultHandler struct {
	svc ports.ITaskResultService
}

// NewTaskResultHandler creates a new TaskResultHandler instance
func NewTaskResultHandler(svc ports.ITaskResultService, rg *gin.RouterGroup) TaskResultHandler {
	taskRouteGroup := rg.Group("/task-result")
	handler := TaskResultHandler{
		svc,
	}

	taskRouteGroup.POST("/", handler.SubmitTaskResult)
	taskRouteGroup.GET("/", handler.ListTaskResults)
	taskRouteGroup.GET("/:id", handler.GetTaskResult)
	taskRouteGroup.PUT("/", handler.UpdateTaskResult)
	taskRouteGroup.DELETE("/:id", handler.DeleteTaskResult)

	return handler
}

// submitRequest represents the request body for creating a task result
type submitRequest struct {
	Name    string  `json:"name" binding:"required" example:"John Doe"`
	Score   float64 `json:"score" binding:"required" example:"6.5"`
	Comment string  `json:"comment" binding:"required" example:"This is a comment for submitted task"`
}

// taskResultResponse represents a task result response body
type taskResultResponse struct {
	ID      string  `json:"id" example:"aaa-bbb-ccc-ddd"`
	Name    string  `json:"name" example:"John Doe"`
	Score   float64 `json:"score" example:"6.5"`
	Comment string  `json:"comment" example:"This is a comment for submitted task"`
}

// newUserResponse is a helper function to create a response body for handling user data
func newTaskResultResponse(t *taskResultDomain.TaskResultEntity) taskResultResponse {
	return taskResultResponse{
		ID:    t.ID,
		Name:  t.Name,
		Score: float64(t.Score),
	}
}

func (h TaskResultHandler) SubmitTaskResult(ctx *gin.Context) {
	var req submitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	taskResult := taskResultDomain.TaskResultEntity{
		ID:      uuid.NewString(),
		Name:    req.Name,
		Score:   req.Score,
		Comment: req.Comment,
	}

	_, err := h.svc.SubmitTask(ctx, &taskResult)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newTaskResultResponse(&taskResult)
	handleSuccess(ctx, rsp)
}

// listUsersRequest represents the request body for listing users
type listUsersRequest struct {
	Skip  uint64 `form:"skip" binding:"min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

func (h TaskResultHandler) ListTaskResults(ctx *gin.Context) {
	var req listUsersRequest
	var taskList []taskResultResponse
	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	tasks, err := h.svc.ListTaskResults(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, t := range tasks {
		taskList = append(taskList, newTaskResultResponse(&t))
	}

	total := uint64(len(taskList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, taskList, "taskResults")
	handleSuccess(ctx, rsp)
}

// getTaskResultRequest represents the request body for getting a task result
type getTaskResultRequest struct {
	ID string `uri:"id" binding:"required" example:"1"`
}

func (h TaskResultHandler) GetTaskResult(ctx *gin.Context) {
	var req getTaskResultRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	taskResult, err := h.svc.GetTaskResult(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newTaskResultResponse(taskResult)
	handleSuccess(ctx, rsp)
}

// updateTaskResultRequest represents the request body for updating a task result
type updateTaskResultRequest struct {
	Name    string  `json:"name" binding:"omitempty,required" example:"John Doe"`
	Email   string  `json:"email" binding:"omitempty,required,email" example:"test@example.com"`
	Score   float64 `json:"score" binding:"omitempty,required" example:"6.5"`
	Comment string  `json:"comment" binding:"omitempty,required" example:"This is a comment for submitted task"`
}

func (h TaskResultHandler) UpdateTaskResult(ctx *gin.Context) {
	var req updateTaskResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	task := taskResultDomain.TaskResultEntity{
		ID:      idStr,
		Name:    req.Name,
		Score:   req.Score,
		Comment: req.Comment,
	}

	_, err := h.svc.UpdateTaskResult(ctx, &task)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newTaskResultResponse(&task)
	handleSuccess(ctx, rsp)
}

// deleteTaskResultRequest represents the request body for deleting a task result
type deleteTaskResultRequest struct {
	ID string `uri:"id" binding:"required" example:"1"`
}

func (h TaskResultHandler) DeleteTaskResult(ctx *gin.Context) {
	var req deleteTaskResultRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := h.svc.DeleteTaskResult(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
