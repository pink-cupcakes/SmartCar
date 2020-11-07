package loghelper

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"app_api/shared"
	testhelper "app_api/shared/testhelpers"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

var hook *test.Hook = test.NewGlobal()

func TestAssignRequestID(t *testing.T) {
	ctx := testhelper.GetTestContext()
	c := AssignRequestID(ctx)
	d := AssignRequestID(c)

	assert.NotNil(t, c.Value(ContextKeyRequestID), "The request ID should not be nil.")
	assert.NotNil(t, d.Value(ContextKeyRequestID), "The request ID should not be nil.")

	assert.NotEqual(t, c.Value(ContextKeyRequestID), d.Value(ContextKeyRequestID), "The request ID should have changed from the original value.")
}

func TestGetRequestID(t *testing.T) {
	ctx := testhelper.GetTestContext()
	c := AssignRequestID(ctx)
	d := AssignRequestID(c)

	cg := GetRequestID(c)
	dg := GetRequestID(d)

	assert.Equal(t, c.Value(ContextKeyRequestID), cg, "The request IDs should be the same.")
	assert.Equal(t, d.Value(ContextKeyRequestID), dg, "The request IDs should be the same.")
}

func TestLogErrors(t *testing.T) {
	expectedError := "Bad Request"
	err := shared.NewAPIError(http.StatusBadRequest, errors.New(expectedError), "Internal Error").
		SetInternalErrorMessage(expectedError)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyRequestID, "Test-Request-ID")
	reqID := ctx.Value(ContextKeyRequestID)

	LogErrors(ctx, err)

	assert.Equal(t, 1, len(hook.Entries), "The logged error should exist.")

	actualError := hook.LastEntry().Data["Error"]
	assert.Equal(t, expectedError, actualError, "The logged raw error message should be the same as the error message.")

	actualError = hook.LastEntry().Data["InternalError"]
	assert.Equal(t, expectedError, actualError, "The logged error message should be the same as the error message.")

	expectedErrorCode := http.StatusBadRequest
	actualErrorCode := hook.LastEntry().Data["ErrorCode"]
	assert.Equal(t, expectedErrorCode, actualErrorCode, "The logged error code should be the same as the error message.")

	expectedError = "Internal Error"
	actualError = hook.LastEntry().Data["ClientError"]
	assert.Equal(t, expectedError, actualError, "The logged client error message should be the same as the error message.")

	actualID := hook.LastEntry().Data["RequestID"]
	assert.Equal(t, reqID, actualID, "The logged error requestID should be the same as the error message.")

	hook.Reset()
	assert.Equal(t, 0, len(hook.Entries), "The mocked logger should be reset.")
}

func TestLogErrorsNoCTX(t *testing.T) {
	expectedError := "Not Found"
	err := shared.NewAPIError(http.StatusBadRequest, errors.New(expectedError), "Internal Error").
		SetInternalErrorMessage(expectedError)
	LogErrorsNoCTX(err)

	assert.Equal(t, 1, len(hook.Entries), "The logged error should exist.")

	actualError := hook.LastEntry().Data["Error"]
	assert.Equal(t, expectedError, actualError, "The logged error should be the same as the error message.")

	actualError = hook.LastEntry().Data["InternalError"]
	assert.Equal(t, expectedError, actualError, "The logged internal error message should be the same as the error message.")

	expectedErrorCode := http.StatusBadRequest
	actualErrorCode := hook.LastEntry().Data["ErrorCode"]
	assert.Equal(t, expectedErrorCode, actualErrorCode, "The logged error code should be the same as the error message.")

	expectedError = "Internal Error"
	actualError = hook.LastEntry().Data["ClientError"]
	assert.Equal(t, expectedError, actualError, "The logged client error message should be the same as the error message.")

	hook.Reset()
	assert.Equal(t, 0, len(hook.Entries), "The mocked logger should be reset.")
}
