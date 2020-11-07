package shared

import (
	"fmt"
	"net/http"
	"runtime"
)

// APIError ... error wrapper for more RESTful errors
type APIError struct {
	ErrorCode              int
	ErrorMessage           error
	InternalErrorMessage   string
	ClientErrorMessage     string
	ValidationErrorMessage string
	file                   string
	line                   int
	funcName               string
	withCallerInf          bool
}

// NewAPIError return error wrapper
// this method will return APIError even if err = 'nil'
func NewAPIError(code int, err error, clientMessage string) *APIError {
	apiError := &APIError{
		ErrorCode:            code,
		ErrorMessage:         err,
		InternalErrorMessage: clientMessage,
		ClientErrorMessage:   clientMessage,
	}
	apiError.validateStatusCode()
	apiError.caller()
	return apiError
}

func (e *APIError) SetInternalErrorMessage(msg string) *APIError {
	e.InternalErrorMessage = msg
	return e
}

func (e *APIError) SetClientErrorMessage(msg string) *APIError {
	e.ClientErrorMessage = msg
	return e
}

func (e *APIError) SetValidationErrorMessage(msg string) *APIError {
	e.ValidationErrorMessage = msg
	return e
}

func (e *APIError) Caller() string {
	if !e.withCallerInf {
		return ""
	}
	return fmt.Sprintf("Called from %s; %s#%d", e.funcName, e.file, e.line)
}

func (e *APIError) caller() {
	var pc uintptr
	pc, e.file, e.line, e.withCallerInf = runtime.Caller(2)
	if e.withCallerInf {
		e.funcName = "undefined"
		if details := runtime.FuncForPC(pc); details != nil {
			e.funcName = details.Name()
		}
	}
	return
}

// validateStatusCode prevent a panic if status code will be wrong
func (e *APIError) validateStatusCode() {
	if e.ErrorCode < 100 || e.ErrorCode > 999 {
		e.ErrorCode = http.StatusInternalServerError
	}
}
