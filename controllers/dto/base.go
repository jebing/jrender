package dto

import "fmt"

type ErrorTitle string

// make the title of the error standardized
const (
	ErrorTitleDeviceConnection   ErrorTitle = "Device Connection Issue"
	ErrorTitleConfig             ErrorTitle = "Config Issue"
	ErrorTitleTimeout            ErrorTitle = "Timeout"
	ErrorTitleMissort            ErrorTitle = "Missort"
	ErrorTitleExternalService    ErrorTitle = "Error from External Service"
	ErrorTitleUserInput          ErrorTitle = "User Error"
	ErrorTitleUnknownParcel      ErrorTitle = "Unknown Parcel"
	ErrorTitleInvalidParcelState ErrorTitle = "Invalid Parcel State"
	ErrorTitleInternalServer     ErrorTitle = "Internal Server Error"
)

type Error struct {
	// This field is not exported as a JSON value; this is purely for internal use to allow the definition of the status code that should be returned
	Status       int        `json:"status,omitempty"`
	Title        ErrorTitle `json:"title,omitempty"`
	ErrorMessage string     `json:"error,omitempty"`
	ErrorCode    int        `json:"error_code,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v - %v", e.Title, e.ErrorMessage)
}

func NewError(status int, title ErrorTitle, message string) *Error {
	return &Error{
		Status:       status,
		Title:        title,
		ErrorMessage: message,
	}
}
