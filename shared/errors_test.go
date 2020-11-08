package shared

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIErrorSuccess(t *testing.T) {
	err := errors.New("This is a test error")
	clientErrorMessage := ""

	successAPIError := &APIError{
		ErrorCode:          http.StatusInternalServerError,
		ErrorMessage:       err,
		ClientErrorMessage: clientErrorMessage,
	}

	apiError := NewAPIError(http.StatusInternalServerError, err, clientErrorMessage)
	assert.Equal(t, successAPIError.ErrorCode, apiError.ErrorCode)
	assert.Equal(t, successAPIError.ErrorMessage, apiError.ErrorMessage)
	assert.Equal(t, successAPIError.ClientErrorMessage, apiError.ClientErrorMessage)
}

func TestSetInternalErrorMessageSuccess(t *testing.T) {
	err := errors.New("This is a test error")
	clientErrorMessage := "Client error message test"
	internalErrorMessage := "This is an internal message test"

	successAPIError := &APIError{
		ErrorCode:          http.StatusInternalServerError,
		ErrorMessage:       err,
		ClientErrorMessage: clientErrorMessage,
	}
	successAPIError.SetInternalErrorMessage(internalErrorMessage)

	assert.Equal(t, successAPIError.InternalErrorMessage, internalErrorMessage)
}
