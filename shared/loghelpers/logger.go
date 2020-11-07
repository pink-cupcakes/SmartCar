package loghelper

import (
	"context"
	"fmt"
	"net/http"

	"app_api/shared"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type key string

// ContextKeyRequestID ... key for the request ID in context
const ContextKeyRequestID key = "requestID"

// ContextKeyRequestPath ... key for the request path in context
const ContextKeyRequestPath key = "requestPath"

// AssignRequestID ... create new request ID and add it to the context
func AssignRequestID(ctx context.Context) context.Context {
	reqID := uuid.New()
	return context.WithValue(ctx, ContextKeyRequestID, reqID.String())
}

// GetRequestID ... get request ID from context
func GetRequestID(ctx context.Context) string {
	reqID := ctx.Value(ContextKeyRequestID)
	if ret, ok := reqID.(string); ok {
		return ret
	}
	return ""
}

// AssignRequestPath ... add requested path and query to the context
func AssignRequestPath(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, ContextKeyRequestPath, fmt.Sprintf("%s.%s", r.Method, r.RequestURI))
}

// GetRequestPath ... get request ID from context
func GetRequestPath(ctx context.Context) string {
	reqPath := ctx.Value(ContextKeyRequestPath)
	if ret, ok := reqPath.(string); ok {
		return ret
	}
	return ""
}

// LogErrors ... helper function to log errors
func LogErrors(ctx context.Context, err *shared.APIError) {
	if err == nil {
		return
	}
	var errMsg string
	if err.ErrorMessage != nil {
		errMsg = err.ErrorMessage.Error()
	}
	logMessage := log.WithContext(ctx).WithFields(log.Fields{
		"Error":         errMsg,
		"ErrorCode":     err.ErrorCode,
		"InternalError": err.InternalErrorMessage,
		"ClientError":   err.ClientErrorMessage,
		"RequestID":     GetRequestID(ctx),
		"Request":       GetRequestPath(ctx),
	})
	if err.ValidationErrorMessage != "" {
		logMessage = logMessage.WithFields(log.Fields{"ValidationError": err.ValidationErrorMessage})
	}
	caller := err.Caller()
	if caller != "" {
		logMessage = logMessage.WithFields(log.Fields{"Caller": caller})
	}
	logMessage.Error()
}

// LogErrors ... helper function to log custom messages
func LogCustomMessage(ctx context.Context, requestID, logLevel, action, message string) {
	logMessage := log.WithContext(ctx).WithFields(log.Fields{
		"Action":    action,
		"RequestID": requestID,
	})

	switch logLevel {
	case "debug":
		logMessage.Debug(message)
	case "info":
		logMessage.Info(message)
	case "warn":
		logMessage.Warn(message)
	case "error":
		logMessage.Error(message)
	}
}

// LogErrorsNoCTX ... helper function to log errors without context
func LogErrorsNoCTX(err *shared.APIError) {
	if err == nil {
		return
	}
	var errMsg string
	if err.ErrorMessage != nil {
		errMsg = err.ErrorMessage.Error()
	}
	log.WithFields(log.Fields{
		"Error":           errMsg,
		"ErrorCode":       err.ErrorCode,
		"InternalError":   err.InternalErrorMessage,
		"ClientError":     err.ClientErrorMessage,
		"ValidationError": err.ValidationErrorMessage,
		"Caller":          err.Caller(),
	}).Error()
}
