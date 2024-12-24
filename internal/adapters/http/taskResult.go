package http

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	taskResultDomain "github.com/lk153/quizgame-ai-serving/internal/core/domains/taskResult"
	"github.com/lk153/quizgame-ai-serving/internal/core/ports"
	copilotagent "github.com/lk153/quizgame-ai-serving/lib/copilotAgent"
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
	taskRouteGroup.POST("/assess", handler.AssessIELTS)
	taskRouteGroup.POST("/upload", handler.Uploadfile)

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
func newTaskResultResponse(t *taskResultDomain.TaskResultEntity) *taskResultResponse {
	if t == nil {
		return nil
	}

	return &taskResultResponse{
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
	var taskListResp []*taskResultResponse
	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	taskResults, err := h.svc.ListTaskResults(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, t := range taskResults {
		taskListResp = append(taskListResp, newTaskResultResponse(&t))
	}

	total := uint64(len(taskListResp))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, taskListResp, "taskResults")
	handleSuccess(ctx, rsp)
}

// getTaskResultRequest represents the request body for getting a task result
type getTaskResultRequest struct {
	ID string `uri:"id" binding:"required" example:"4bf0b061-3926-425f-af89-7b4edb1db389"`
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

// assessRequest represents the request body for IELTS Writing Task assessment
type assessRequest struct {
	TaskType        uint8  `json:"task_type" binding:"required" example:"This is a writing type"`
	TaskRequirement string `json:"task_requirement" binding:"required" example:"This is a writing task"`
	TaskFile        string `json:"task_file"`
	CandidateText   string `json:"candidate_text" binding:"required" example:"This is a candidate text"`
}

func (h TaskResultHandler) AssessIELTS(ctx *gin.Context) {
	var req assessRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	// result, err := copilotagent.DoAssessment(ctx, 123, copilotagent.InputTask{
	// 	TaskType:        req.TaskType,
	// 	TaskRequirement: req.TaskRequirement,
	// 	CandidateText:   req.CandidateText,
	// })
	result, err := copilotagent.DoAssessmentV1(ctx, "4a546073-385a-4e92-8e9a-abbe51268cc5", copilotagent.InputTask{
		TaskType:        req.TaskType,
		TaskRequirement: req.TaskRequirement,
		CandidateText:   req.CandidateText,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, result)
}

func (h TaskResultHandler) Uploadfile(ctx *gin.Context) {
	// single file
	uploadedFile, _ := ctx.FormFile("file")
	log.Println(uploadedFile.Filename)

	tempFile, err := os.CreateTemp("tmp", "prefix")
	if err != nil {
		err = errors.New(fmt.Sprintf("create temp file err: %s", err.Error()))
		handleError(ctx, err)
		return
	}

	defer os.Remove(tempFile.Name())
	ext := filepath.Ext(uploadedFile.Filename)

	// Upload the file to specific dst.
	filename := filepath.Base(tempFile.Name())
	if err := ctx.SaveUploadedFile(uploadedFile, filename+ext); err != nil {
		err = errors.New(fmt.Sprintf("Upload file err: %s", err.Error()))
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
