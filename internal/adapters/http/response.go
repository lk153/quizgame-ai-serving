package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	domainErr "github.com/lk153/quizgame-ai-serving/internal/core/domains/error"
	errLib "github.com/lk153/quizgame-ai-serving/lib/errors"
)

// response represents a response body format
type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}

var errHTTPStatuses = map[error]int{
	domainErr.ErrInternal:                   http.StatusInternalServerError,
	domainErr.ErrDataNotFound:               http.StatusNotFound,
	domainErr.ErrConflictingData:            http.StatusConflict,
	domainErr.ErrInvalidCredentials:         http.StatusUnauthorized,
	domainErr.ErrUnauthorized:               http.StatusUnauthorized,
	domainErr.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domainErr.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domainErr.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domainErr.ErrInvalidToken:               http.StatusUnauthorized,
	domainErr.ErrExpiredToken:               http.StatusUnauthorized,
	domainErr.ErrForbidden:                  http.StatusForbidden,
	domainErr.ErrNoUpdatedData:              http.StatusBadRequest,
}

func handleError(ctx *gin.Context, err error) {
	statusCode, ok := errHTTPStatuses[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := errLib.ParseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

// validationError sends an error response for some specific request validation error
func validationError(ctx *gin.Context, err error) {
	errMsgs := errLib.ParseError(err)
	errRsp := newErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}
