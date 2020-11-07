package testhelper

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TEST HELPER FUNCTIONS

// CheckResponseCode ... assert equal for http response codes
func CheckResponseCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, "expected "+strconv.Itoa(expected)+" but got "+strconv.Itoa(actual))
}

// GetTestContext ... get fake context for testing
func GetTestContext() context.Context {
	// set up context, simulating request.Context
	type key string
	const (
		reqID key = "requestID"
	)
	ctx := context.Background()
	ctx = context.WithValue(ctx, reqID, "Test-Request-ID")
	return ctx
}
