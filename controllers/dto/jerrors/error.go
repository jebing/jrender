package jerrors

import (
	"encoding/json"
	"net/http"
)

type ErrorResp struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error ErrorResp `json:"error"`
}

func NewErrorResp(code int, message string) ErrorResp {
	return ErrorResp{
		Code:    code,
		Message: message,
	}
}

func (e ErrorResp) Error() string {
	return e.Message
}

// Common error constructors
func BadRequest(message string) ErrorResp {
	return NewErrorResp(http.StatusBadRequest, message)
}

func Unauthorized(message string) ErrorResp {
	return NewErrorResp(http.StatusUnauthorized, message)
}

func NotFound(message string) ErrorResp {
	return NewErrorResp(http.StatusNotFound, message)
}

func Conflict(message string) ErrorResp {
	return NewErrorResp(http.StatusConflict, message)
}

func UnprocessableEntity(message string) ErrorResp {
	return NewErrorResp(http.StatusUnprocessableEntity, message)
}

func TooManyRequests(message string) ErrorResp {
	return NewErrorResp(http.StatusTooManyRequests, message)
}

func InternalServerError(message string) ErrorResp {
	return NewErrorResp(http.StatusInternalServerError, message)
}

// WriteErrorResponse writes an error response to the HTTP response writer
func WriteErrorResponse(w http.ResponseWriter, err ErrorResp) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	
	response := ErrorResponse{Error: err}
	json.NewEncoder(w).Encode(response)
}
