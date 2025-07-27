package dto

import "revonoir.com/jbilling/controllers/dto/jerrors"

type Response[T interface{}] struct {
	Data  T                  `json:"data,omitempty"`
	Error *jerrors.ErrorResp `json:"error,omitempty"`
}

func NewErrorResponse(code int, message string) Response[any] {
	resp := Response[any]{
		Data: nil,
		Error: &jerrors.ErrorResp{
			Code:    code,
			Message: message,
		},
	}
	return resp
}
